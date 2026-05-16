package s3

import (
	"context"
	"errors"
)

// FileStorage is a stub implementation for S3.
type FileStorage struct{}

// NewFileStorage creates a new S3 storage adapter.
func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

// Put uploads a file to storage.
func (s *FileStorage) Put(ctx context.Context, key string, data []byte, contentType string) error {
	_ = ctx
	_ = key
	_ = data
	_ = contentType
	return errors.New("not implemented")
}

// GetURL returns a public URL to the stored file.
func (s *FileStorage) GetURL(ctx context.Context, key string) (string, error) {
	_ = ctx
	_ = key
	return "", errors.New("not implemented")
}
