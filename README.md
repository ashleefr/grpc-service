У меня сейчас такая архитектура:
grpc-file-service
├── README.md
├── cmd
│        └── server
│            └── start.go
├── go.mod
├── go.sum
├── internal
│        ├── server
│        │   └── server.go
│        └── storage
│            └── storage.go
├── main.go
├── proto
│        ├── file_service.pb.go
│        ├── file_service.proto
│        └── file_service_grpc.pb.go
└── storage

start.go
```
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
		log.Fatalf("Failed to start the gRPC server: %v", err)
	}
}
```

server.go
```
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

func NewFileServiceServer(storage *storage.FileStorage) pb.FileServiceServer {
	return &fileServiceServer{
		fileStorage:           storage,
		uploadDownloadLimiter: make(chan struct{}, 10),  // Ограничение на 10 одновременных запросов загрузки/скачивания
		listLimiter:           make(chan struct{}, 100), // Ограничение на 100 одновременных запросов списка
	}
}

func (s *fileServiceServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	// Acquire limiter slot to restrict simultaneous uploads/downloads
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }() // Release limiter slot on function exit

	// Delegate the upload process to FileStorage
	return s.fileStorage.Upload(stream)
}

func (s *fileServiceServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// Acquire limiter slot to restrict simultaneous list requests
	s.listLimiter <- struct{}{}
	defer func() { <-s.listLimiter }() // Release limiter slot on function exit

	// Delegate the list retrieval process to FileStorage
	return s.fileStorage.List()
}

func (s *fileServiceServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {
	// Acquire limiter slot to restrict simultaneous uploads/downloads
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }() // Release limiter slot on function exit

	// Delegate the download process to FileStorage
	return s.fileStorage.Download(req.Filename, stream)
}
```

storage.go
```
package storage

import (
	"errors"
	pb "grpc-file-service/proto"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	storagePath string
}

func NewFileStorage(path string) (*FileStorage, error) {
	// Check if a directory already exists
	if _, err := os.Stat(path); err != nil {
		// Create a directory
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &FileStorage{storagePath: path}, nil
}

func (s *FileStorage) Upload(stream pb.FileService_UploadFileServer) error {
	// Open a file for writing
	var filename string
	var file *os.File

	// Read a file from the stream
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// If the stream is closed, break the loop
			break
		}
		if err != nil {
			// If an error occurs, return the error
			return err
		}

		// If the file is not opened yet, open it
		if file == nil {
			// Create a full path to the file
			filename = filepath.Join(s.storagePath, req.GetFilename())
			// Create the file
			file, err = os.Create(filename)
			if err != nil {
				// If an error occurs, return the error
				return err
			}
			// Close the file when the function is finished
			defer file.Close()
		}

		// Write the data to the file
		if _, err := file.Write(req.GetData()); err != nil {
			// If an error occurs, return the error
			return err
		}
	}

	// Send a response to the client
	return stream.SendAndClose(&pb.UploadFileResponse{Message: "Файл успешно загружен"})
}

func (s *FileStorage) List() (*pb.ListFilesResponse, error) {
	files, err := os.ReadDir(s.storagePath)
	if err != nil {
		return nil, err
	}

	var fileInfos []*pb.FileInfo
	for _, f := range files {
		if f.IsDir() {
			// Skip directories
			continue
		}

		filePath := filepath.Join(s.storagePath, f.Name())
		fileStat, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, &pb.FileInfo{
			Filename:  f.Name(),
			CreatedAt: fileStat.ModTime().Format(time.RFC3339),
			UpdatedAt: fileStat.ModTime().Format(time.RFC3339),
		})
	}

	return &pb.ListFilesResponse{Files: fileInfos}, nil
}

func (s *FileStorage) Download(filename string, stream pb.FileService_DownloadFileServer) error {
	// Create a full path to the file
	filePath := filepath.Join(s.storagePath, filename)
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If an error occurs, return an error with a message
		return errors.New("файл не найден")
	}
	// Close the file when the function is finished
	defer file.Close()

	// Create a buffer for reading the file
	buf := make([]byte, 1024)
	for {
		// Read from the file
		n, err := file.Read(buf)
		// If the end of the file is reached, break the loop
		if err == io.EOF {
			break
		}
		// If an error occurs, return the error
		if err != nil {
			return err
		}

		// Send the data to the client
		if err := stream.Send(&pb.DownloadFileResponse{Data: buf[:n]}); err != nil {
			// If an error occurs, return the error
			return err
		}
	}
	return nil
}
```

file_service.proto
```
syntax = "proto3";

package proto;

option go_package = "/proto;proto";

service FileService {
  // Upload a file to the server.
  //
  // The request message stream is intended to be a single message with a
  // filename and the file contents. The response message is a single message
  // with a blank filename and a status code.
  rpc UploadFile (stream UploadFileRequest) returns (UploadFileResponse);
  // List the files stored on the server.
  //
  // The request message is empty, and the response message is a stream of
  // FileInfo messages.
  rpc ListFiles (ListFilesRequest) returns (ListFilesResponse);
  // Download a file from the server.
  //
  // The request message contains the filename of the desired file.
  // The response message is a stream of data chunks of the file.
  rpc DownloadFile (DownloadFileRequest) returns (stream DownloadFileResponse);
}

message UploadFileRequest {
  string filename = 1;
  bytes data = 2;
}

message UploadFileResponse {
  string message = 1;
}

message ListFilesRequest {}

message ListFilesResponse {
  repeated FileInfo files = 1;
}

message FileInfo {
  string filename = 1;
  string created_at = 2;
  string updated_at = 3;
}

message DownloadFileRequest {
  string filename = 1;
}

message DownloadFileResponse {
  bytes data = 1;
}
```

main.go
```
package main

import "grpc-file-service/cmd/server"

func main() {
	// Run the gRPC server
	server.Run()
}

```


Измени ее под архитектуру, что будет ниже, с интерфейсами, структурами и прочим. Также нужно добавить связь с БД PostgreSQL

Вот еще мои заметки, когда проверяли изначальный проект:
"
- Бд, хранящая данные о файлах
- Сам сервис с бд не должен общаться с сервисом логики хранения файлов
- Все через юзкейсы

- Допустим, в юзкейсе вызывается метод сохранения картинки, метод сохранения возвращает в юзкейс данные о сохраненном файле
- Эти данные из юзкейса передаются в метод сохранения данных в БД

- Мы заставляем слои общаться через модель
- Создать модели, они должны храниться в отдельном файле (известны всем). Для общения между юзкейсами и сервисами, а второй - между юзкейсами и транспортом. 1 - модели сервисов, 2 - модель юзкейсов

- Отдельный слой транспорта - вызов gRPC сервера

- Никто ничего не знает о других реализациях из других слоев, все описывается через интерфейсы и коннектится в старте
"