package server

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc-file-service/internal/server"
	"grpc-file-service/internal/storage"
	pb "grpc-file-service/proto"
	"log"
	"net"
)

var (
	port = flag.Int("port", 1337, "The server port")
)

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}

	fileStorage, err := storage.NewFileStorage("storage")
	if err != nil {
		log.Fatalf("Failed to initialize the file storage: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, server.NewFileServiceServer(fileStorage))

	log.Printf("Server listening at %v\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Не удалось запустить gRPC сервер: %v", err)
	}
}
