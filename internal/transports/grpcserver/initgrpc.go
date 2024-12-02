package grpcserver

import (
	"context"
	"log"
	"net"

	"grpc-file-service/internal/models/modelssvc"
	"grpc-file-service/internal/usecases/fileusecase"
	pb "grpc-file-service/pkg/proto"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	usecase *fileusecase.UseCase
	pb.UnimplementedFileServiceServer
}

func NewGRPCServer(uc *fileusecase.UseCase) *GRPCServer {
	return &GRPCServer{usecase: uc}
}

func (s *GRPCServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	file := &modelssvc.File{
		ID:   "some_generated_id",
		Name: req.GetName(),
		Data: req.GetData(),
	}
	err := s.usecase.UploadFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &pb.UploadFileResponse{Message: "File uploaded successfully"}, nil
}

func (s *GRPCServer) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	file, err := s.usecase.GetFile(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetFileResponse{
		Name: file.Name,
		Data: file.Data,
	}, nil
}

func InitGRPCServer(port string, uc *fileusecase.UseCase) {
	server := grpc.NewServer()
	grpcServer := NewGRPCServer(uc)
	pb.RegisterFileServiceServer(server, grpcServer)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("gRPC server started on port %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
