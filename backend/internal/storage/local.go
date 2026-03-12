package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	baseDir string
}

func NewLocal(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir}
}

func (l *LocalStorage) Save(_ context.Context, key string, body io.Reader, _ int64, _ string) error {
	path := filepath.Join(l.baseDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (l *LocalStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(l.baseDir, filepath.FromSlash(key))
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return file, nil
}
