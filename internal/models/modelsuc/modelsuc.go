package modelsuc

import "time"

// FileInfo представляет информацию о файле
type FileInfo struct {
	Filename   string    // Имя файла
	Size       int64     // Размер файла в байтах
	UploadDate time.Time // Дата загрузки файла
}
