package modelssvc

import (
	"grpc-file-service/internal/models/modelsuc"
	"time"
)

// FileServiceModel — модель для представления файла в сервисе
type FileServiceModel struct {
	ID         int       // Уникальный идентификатор файла в базе данных
	Filename   string    // Имя файла
	Size       int64     // Размер файла в байтах
	UploadDate time.Time // Дата загрузки файла
}

// ToUseCaseModel конвертирует FileServiceModel в FileInfo из UseCase слоя
func (fsm *FileServiceModel) ToUseCaseModel() *modelsuc.FileInfo {
	return &modelsuc.FileInfo{
		Filename:   fsm.Filename,
		Size:       fsm.Size,
		UploadDate: fsm.UploadDate,
	}
}
