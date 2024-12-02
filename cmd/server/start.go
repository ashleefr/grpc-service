package main

import (
	"grpc-file-service/internal/services/db"
	"grpc-file-service/internal/transports/grpcserver"
	"grpc-file-service/internal/usecases/fileusecase"
)

func main() {
	repo := db.NewDB()
	usecase := fileusecase.New(repo)
	grpcserver.InitGRPCServer("50051", usecase)
}
