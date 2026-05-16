package serviceInterface

import "context"

// FileStorage stores uploaded files such as proof of payment.
type FileStorage interface {
	Put(ctx context.Context, key string, data []byte, contentType string) error
	GetURL(ctx context.Context, key string) (string, error)
}
