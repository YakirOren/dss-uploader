package models

import (
	"io/fs"
	"os"
	"time"
)

type FileMetadata struct {
	FileName    string `bson:"name"`
	FileSize    int64  `bson:"size"`
	IsDirectory bool   `bson:"isDirectory"`
	Fragments   []Fragment
}

func NewFileMetadata(name string, size int64, isDir bool, fragments []Fragment) *FileMetadata {
	return &FileMetadata{
		FileName:    name,
		FileSize:    size,
		IsDirectory: isDir,
		Fragments:   fragments,
	}
}

func (f FileMetadata) Name() string {
	return f.FileName
}

func (f FileMetadata) Size() int64 {
	return f.FileSize
}

func (f FileMetadata) IsDir() bool {
	return f.IsDirectory
}

func (f FileMetadata) Mode() fs.FileMode {
	return os.ModePerm
}

func (f FileMetadata) ModTime() time.Time {
	return time.Now()
}

func (f FileMetadata) Sys() interface{} {
	return nil
}
