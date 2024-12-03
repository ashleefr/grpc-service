package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-file-service/internal/usecases/fileusecase"
	"grpc-file-service/pkg/proto"
	"log"
	"net"
)

type GRPCServer struct {
	useCase fileusecase.FileUseCase
	proto.UnimplementedFileServiceServer
}

func NewGRPCServer(useCase *fileusecase.FileUseCase) *GRPCServer {
	return &GRPCServer{
		useCase: useCase,
	}
}

func (s *GRPCServer) Start(port int) error {
	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterFileServiceServer(grpcServer, s)

	log.Printf("gRPC server started on :%d", port)
	return grpcServer.Serve(lis)
}

func (s *GRPCServer) Stop() error {
	// Здесь можно добавить логику для graceful shutdown
	log.Println("Stopping gRPC server...")
	return nil
}
