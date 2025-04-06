package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/KamigamiNoGigan/grpc/pkg/server_api_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	pb.UnimplementedUserAPIServer
	storage   []*pb.Info
	uploadDir string

	uploadSem   chan struct{}
	downloadSem chan struct{}
	infoSem     chan struct{}

	storageMu sync.Mutex
}

func NewServer() *server {
	s := &server{
		uploadSem:   make(chan struct{}, 10),
		downloadSem: make(chan struct{}, 10),
		infoSem:     make(chan struct{}, 100),
		uploadDir:   "storage",
	}

	if _, err := os.Stat("storage.json"); errors.Is(err, os.ErrNotExist) {
		os.Create("storage.json")
	}

	b, _ := os.ReadFile("storage.json")

	err := json.Unmarshal(b, &s.storage)

	if err != nil {
		log.Fatalf("failed to unmarshal storage.json: %v", err)
	}

	return s
}

func (s *server) Upload(stream pb.UserAPI_UploadServer) error {
	s.uploadSem <- struct{}{}
	defer func() {
		<-s.uploadSem
	}()

	s.uploadDir = "storage"
	var (
		file       *os.File
		info       *pb.Info
		filename   string
		firstChunk = true
	)

	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	default:
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				s.storageMu.Lock()
				s.storage = append(s.storage, info)

				data, err := json.Marshal(s.storage)
				if err != nil {
					return err
				}
				err = os.WriteFile("storage.json", data, 0777)
				if err != nil {
					log.Printf("write file error: %v", err)
					s.storageMu.Unlock()
					return err
				}
				s.storageMu.Unlock()
				return stream.SendAndClose(&wrapperspb.BoolValue{Value: true})
			}
			if err != nil {
				if file != nil {
					file.Close()
					os.Remove(file.Name())
				}
				return status.Errorf(codes.Internal, "failed to receive chuck: %v", err)
			}

			if firstChunk {
				firstChunk = false
				if chunk.GetFileName() == "" {
					return status.Errorf(codes.Internal, "first chunk must contain metadata")
				}

				filename = chunk.GetFileName()
				s.storageMu.Lock()
				info = FindName(filename, &s.storage)
				s.storageMu.Unlock()
				filePath := filepath.Join(s.uploadDir, filename)

				file, err = os.Create(filePath)
				if err != nil {
					return status.Errorf(codes.Internal, "failed to create file: %v", err)
				}
				continue
			}

			if chunk.GetChunkData() == nil {
				return status.Error(codes.InvalidArgument, "non-first chunk must contain data")
			}

			if _, err := file.Write(chunk.GetChunkData()); err != nil {
				file.Close()
				os.Remove(file.Name())
				return status.Errorf(codes.Internal, "failed to write chunk: %v", err)
			}

		}
	}

}

func (s *server) GetInfo(ctx context.Context, st *emptypb.Empty) (*pb.FileList, error) {
	s.infoSem <- struct{}{}
	defer func() {
		<-s.infoSem
	}()

	return &pb.FileList{Files: s.storage}, nil
}

func (s *server) Download(name *pb.DownloadRequest, stream pb.UserAPI_DownloadServer) error {
	s.downloadSem <- struct{}{}
	defer func() {
		<-s.downloadSem
	}()
	var find bool = false

	s.uploadDir = "storage"
	for _, v := range s.storage {
		if v.FileName == name.FileName {
			find = true
			filePath := filepath.Join(s.uploadDir, v.FileName)
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			buf := make([]byte, 64*1024)
			for {
				select {
				case <-stream.Context().Done():
					return stream.Context().Err()
				default:
					n, err := file.Read(buf)
					if err == io.EOF {
						return nil
					}
					if err != nil {
						return err
					}
					if err := stream.Send(&pb.DownloadResponse{
						ChunkData: buf[:n],
					}); err != nil {
						return err
					}
				}
			}
		}
	}
	if !find {
		return status.Errorf(codes.NotFound, "file not found: %s", name.FileName)
	}
	return nil
}

func FindName(filename string, storage *[]*pb.Info) *pb.Info {
	createdAt := timestamppb.Now()
	updatedAt := timestamppb.Now()
	for i, v := range *storage {
		if v.GetFileName() == filename {
			createdAt = v.CreatedAt
			updatedAt = timestamppb.Now()
			*storage = append((*storage)[:i], (*storage)[i+1:]...)
			return &pb.Info{FileName: filename, CreatedAt: createdAt, UpdatedAt: updatedAt}
		}
	}
	return &pb.Info{FileName: filename, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxConcurrentStreams(100))
	reflection.Register(s)
	pb.RegisterUserAPIServer(s, NewServer())

	log.Println("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
