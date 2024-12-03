package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"grpc-file-service/internal/services/localstorage"
	"grpc-file-service/internal/transports/grpcserver"
	"grpc-file-service/internal/usecases/fileusecase"
)

type App struct {
	DB          *sql.DB
	FileUseCase fileusecase.FileUseCase
	GRPCServer  *grpcserver.GRPCServer
}

func InitApp(dbConnStr, storagePath string) (*App, error) {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	// Инициализация хранилища файлов
	fileStore, err := localstorage.NewFileStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file storage: %w", err)
	}

	// Создаем FileUseCase
	fileUseCase := fileusecase.NewFileUseCase(db, fileStore)

	// Инициализация gRPC сервера
	grpcServer := grpcserver.NewGRPCServer(fileUseCase)

	return &App{
		DB:          db,
		FileUseCase: fileUseCase,
		GRPCServer:  grpcServer,
	}, nil
}
