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

// Run Starts the gRPC server and listens for incoming requests
func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}

	// Creating a new file storage
	fileStorage, err := storage.NewFileStorage("storage")
	if err != nil {
		log.Fatalf("Failed to initialize the file storage: %v", err)
	}

	// Creating a new gRPC server
	grpcServer := grpc.NewServer()

	// Registering the FileService server
	pb.RegisterFileServiceServer(grpcServer, server.NewFileServiceServer(fileStorage))

	log.Printf("Server listening at %v\n", lis.Addr())

	// Starting the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the gRPC server: %v", err)
	}
}
