package s3storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Storage struct {
	client *minio.Client
	bucket string
}

func New(ctx context.Context, cfg Config) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("initialize storage: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}

	if !exists {
		err := client.MakeBucket(
			ctx,
			cfg.Bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &Storage{client: client, bucket: cfg.Bucket}, nil
}

func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, opts)
	if err != nil {
		return fmt.Errorf("upload object: %w", err)
	}

	return nil
}

func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting object: %w", err)
	}

	if _, err := object.Stat(); err != nil {
		_ = object.Close()
		return nil, fmt.Errorf("stat object: %w", err)
	}

	return object, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}

func (s *Storage) PresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned URL: %w", err)
	}

	return url.String(), nil
}
