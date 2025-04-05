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
	storage   map[string]Information
	uploadDir string
}

type Information struct {
	FileName  string                 `json:"name"`
	CreatedAt *timestamppb.Timestamp `json:"created_at"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at"`
}

func (s *server) Upload(stream pb.UserAPI_UploadServer) error {
	var (
		file       *os.File
		filename   string
		createdAt  *timestamppb.Timestamp
		updatedAt  *timestamppb.Timestamp
		firstChunk = true
		info       Information
	)

	if _, err := os.Stat("storage.json"); errors.Is(err, os.ErrNotExist) {
		os.Create("storage.json")
	}

	b, err := os.ReadFile("storage.json")
	if err != nil {
		return err
	}

	if s.storage == nil {
		s.storage = make(map[string]Information)
	}

	err = json.Unmarshal(b, &s.storage)

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {

			s.storage[info.FileName] = info

			data, err := json.Marshal(s.storage)
			if err != nil {
				return err
			}
			os.WriteFile("storage.json", data, 0777)

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
			if v, ok := s.storage[filename]; !ok {
				createdAt = timestamppb.Now()
				updatedAt = timestamppb.Now()
			} else {
				createdAt = v.CreatedAt
				updatedAt = timestamppb.Now()
			}
			s.uploadDir = "storage"
			info.FileName = filename
			info.CreatedAt = createdAt
			info.UpdatedAt = updatedAt

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

func (s *server) GetInfo(ctx context.Context, st *emptypb.Empty) (*pb.FileList, error) {
	return &pb.FileList{}, nil
}

func (s *server) Download(name *pb.DownloadRequest, stream pb.UserAPI_DownloadServer) error {
	return nil
}

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserAPIServer(s, &server{})

	log.Println("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
