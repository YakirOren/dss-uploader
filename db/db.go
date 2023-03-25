package db

import (
	"DSS-uploader/models"
	"context"
)

type DataStore interface {
	WriteFile(ctx context.Context, file models.FileMetadata) (string, error)
	AppendFragment(ctx context.Context, name string, fragment models.Fragment) error
	GetMetadataByName(ctx context.Context, name string) (*models.FileMetadata, bool)
}
