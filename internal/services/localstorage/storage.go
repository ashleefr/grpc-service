package localstorage

import (
	"fmt"
	"os"
)

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	// Создаем директорию, если она не существует
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileStorage{BasePath: basePath}, nil
}

// Пример метода для сохранения файла
func (fs *FileStorage) SaveFile(filename string, content []byte) error {
	filePath := fmt.Sprintf("%s/%s", fs.BasePath, filename)
	return os.WriteFile(filePath, content, 0644)
}
