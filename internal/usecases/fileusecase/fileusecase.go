package fileusecase

import (
	"fmt"
	"grpc-file-service/internal/services/db"
	"grpc-file-service/internal/services/localstorage"
)

type FileUseCase struct {
	db          *db.DB
	fileStorage *localstorage.FileStorage
}

func NewFileUseCase(db *db.DB, fileStorage *localstorage.FileStorage) *FileUseCase {
	return &FileUseCase{
		db:          db,
		fileStorage: fileStorage,
	}
}

func (uc *FileUseCase) SomeDatabaseOperation() error {
	// Пример использования db.DB
	_, err := uc.db.Exec("SELECT 1") // Выполнение запроса
	if err != nil {
		return fmt.Errorf("database operation failed: %v", err)
	}
	return nil
}
