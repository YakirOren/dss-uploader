package upload

import (
	"context"
)

type Client interface {
	Upload(ctx context.Context, path string, file []byte, fragmentID string) error
}
