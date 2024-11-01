package storage

import (
	pb "grpc-file-service/proto"
)

type FileStorage struct {
	storagePath string
}

// NewFileStorage Creating a FileStorage instance
func NewFileStorage(path string) (*FileStorage, error) {
	return nil, nil
}

// Upload Uploading a file
func (s *FileStorage) Upload(stream pb.FileService_UploadFileServer) error {
	return nil
}

// List Viewing a list of files
func (s *FileStorage) List() (*pb.ListFilesResponse, error) {
	return nil, nil
}

// Download Downloading a file
func (s *FileStorage) Download(filename string, stream pb.FileService_DownloadFileServer) error {
	return nil
}
