package db

import (
	"context"
	"errors"
	"myapp/internal/models/modelssvc"
)

type DB struct {
	storage map[string]*modelssvc.File
}

func NewDB() *DB {
	return &DB{storage: make(map[string]*modelssvc.File)}
}

func (db *DB) Save(ctx context.Context, file *modelssvc.File) error {
	db.storage[file.ID] = file
	return nil
}

func (db *DB) FindByID(ctx context.Context, id string) (*modelssvc.File, error) {
	file, ok := db.storage[id]
	if !ok {
		return nil, errors.New("file not found")
	}
	return file, nil
}
