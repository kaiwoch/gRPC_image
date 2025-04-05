package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/KamigamiNoGigan/grpc/pkg/server_api_v1"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("wrong image path")
	}

	path := os.Args[1]
	arr := strings.Split(path, `\`)
	fileName := arr[len(arr)-1]

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserAPIClient(conn)
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
		log.Fatalf(err.Error())
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
	fmt.Println(reply.Value)
}
