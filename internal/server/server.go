package server

import (
	"context"
	"grpc-file-service/internal/storage"
	pb "grpc-file-service/proto"
)

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
	fileStorage           *storage.FileStorage
	uploadDownloadLimiter chan struct{}
	listLimiter           chan struct{}
}

// NewFileServiceServer Creating a new FileService Server
func NewFileServiceServer(storage *storage.FileStorage) pb.FileServiceServer {
	return &fileServiceServer{
		fileStorage:           storage,
		uploadDownloadLimiter: make(chan struct{}, 10),  // Ограничение на 10 одновременных запросов загрузки/скачивания
		listLimiter:           make(chan struct{}, 100), // Ограничение на 100 одновременных запросов списка
	}
}

// UploadFile Uploading a file
func (s *fileServiceServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }()

	return s.fileStorage.Upload(stream)
}

// ListFiles Viewing a list of files
func (s *fileServiceServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	s.listLimiter <- struct{}{}
	defer func() { <-s.listLimiter }()

	return s.fileStorage.List()
}

// DownloadFile Downloading a file
func (s *fileServiceServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }()

	return s.fileStorage.Download(req.Filename, stream)
}
