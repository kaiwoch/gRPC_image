package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/KamigamiNoGigan/grpc/pkg/server_api_v1"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("wrong os.Args")
	}

	path := os.Args[2:]

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserAPIClient(conn)

	if strings.ToLower(os.Args[1]) == "upload" {
		for _, file := range path {
			ClientUpload(c, file)
		}
	}

	if strings.ToLower(os.Args[1]) == "download" {
		for _, fileName := range path {
			ClientDownload(c, fileName)
		}
	}

}

func ClientUpload(c pb.UserAPIClient, path string) {
	arr := strings.Split(path, `\`)
	fileName := arr[len(arr)-1]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	stream, _ := c.Upload(ctx)

	if err := stream.Send(&pb.UploadRequest{
		Data: &pb.UploadRequest_FileName{
			FileName: fileName,
		},
	},
	); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, fileName, err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	buf := make([]byte, 64*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf(err.Error())
		}

		if err := stream.Send(&pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buf[:n],
			},
		}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", c, buf[:n], err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if reply.Value {
		fmt.Printf("file %s successfully uploaded to server\n", fileName)
	}
}

func ClientDownload(c pb.UserAPIClient, fileName string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	searchStream, err := c.Download(ctx, &pb.DownloadRequest{FileName: fileName})

	if err != nil {
		log.Fatalf("%v.Download(%v) = %v", c, fileName, err)
	} else {
		file, err := os.Create(filepath.Join("download", fileName))
		if err != nil {
			log.Fatalf("failed to create file: %v", err)
		}
		defer file.Close()
		for {
			download, err := searchStream.Recv()
			if err == io.EOF {
				fmt.Printf("file %s successfully downloaded from server\n", file.Name())
				break
			}
			if err != nil {
				file.Close()
				os.Remove(file.Name())
				log.Printf("failed to receive chunk: %v\n", err)
				return
			}
			if _, err := file.Write(download.GetChunkData()); err != nil {
				file.Close()
				os.Remove(file.Name())
				log.Printf("failed to write chunk: %v\n", err)
				return
			}
		}
	}
}
