package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
)

type Storage interface {
	Save(ctx context.Context, key string, body io.Reader, size int64, contentType string) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
}

func NewStorage(cfg config.Config) (Storage, error) {
	switch cfg.StorageBackend {
	case "local", "":
		return NewLocal(cfg.StorageDir), nil
	case "s3":
		return NewS3(cfg)
	default:
		return nil, fmt.Errorf("unknown storage backend: %s", cfg.StorageBackend)
	}
}
