package s3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"myapp/domain"
	"myapp/infrastructure/config"
)

const maxFileSize = 5 << 20 // 5 MB

var allowedMIMEs = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// FileStorage implements the FileStorage interface backed by an S3-compatible store.
type FileStorage struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

// New creates a FileStorage client configured for the given S3-compatible endpoint.
func New(cfg config.Config) (*FileStorage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey, cfg.S3SecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		}
		o.UsePathStyle = true
	})

	return &FileStorage{
		client:   client,
		bucket:   cfg.S3Bucket,
		endpoint: cfg.S3Endpoint,
	}, nil
}

// Put validates the file and uploads it to the configured S3 bucket.
func (fs *FileStorage) Put(ctx context.Context, key string, data []byte, contentType string) error {
	if len(data) > maxFileSize {
		return domain.ErrFileTooLarge
	}

	if !allowedMIMEs[contentType] {
		detected := http.DetectContentType(data)
		if !allowedMIMEs[detected] {
			return domain.ErrInvalidMIMEType
		}
		contentType = detected
	}

	_, err := fs.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(fs.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	return err
}

// GetURL returns a public URL for the stored file.
func (fs *FileStorage) GetURL(_ context.Context, key string) (string, error) {
	if fs.endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", fs.endpoint, fs.bucket, key), nil
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", fs.bucket, key), nil
}
