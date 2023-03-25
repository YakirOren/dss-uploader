package db

import (
	"DSS-uploader/models"
	"context"
	"os"
)

type MemoryDataStore struct {
	storage map[string]models.FileMetadata
}

func NewMemoryDataStore() (*MemoryDataStore, error) {
	db := &MemoryDataStore{}
	db.storage = make(map[string]models.FileMetadata)

	return db, nil
}

func (db *MemoryDataStore) WriteFile(_ context.Context, file models.FileMetadata) (string, error) {
	db.storage[file.FileName] = file
	return file.FileName, nil
}

func (db *MemoryDataStore) AppendFragment(_ context.Context, name string, fragment models.Fragment) error {
	file, found := db.storage[name]
	if !found {
		return os.ErrNotExist
	}
	file.Fragments = append(file.Fragments, fragment)

	db.storage[name] = file

	return nil
}
