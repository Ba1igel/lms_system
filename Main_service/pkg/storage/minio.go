package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type minIOStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(ctx context.Context, cfg MinIOConfig) (FileStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &minIOStorage{client: client, bucket: cfg.Bucket}, nil
}

func (s *minIOStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, input.ObjectName, input.Reader, input.Size, minio.PutObjectOptions{
		ContentType: input.ContentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("minio://%s/%s", s.bucket, input.ObjectName), nil
}

func (s *minIOStorage) Download(ctx context.Context, objectURL string) (*DownloadResult, error) {
	objectName, err := objectNameFromURL(objectURL, s.bucket)
	if err != nil {
		return nil, err
	}

	object, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	info, err := object.Stat()
	if err != nil {
		_ = object.Close()
		return nil, err
	}

	return &DownloadResult{
		Reader:      object,
		Size:        info.Size,
		ContentType: info.ContentType,
	}, nil
}

func objectNameFromURL(rawURL, bucket string) (string, error) {
	const minioScheme = "minio://"
	if strings.HasPrefix(rawURL, minioScheme) {
		withoutScheme := strings.TrimPrefix(rawURL, minioScheme)
		prefix := bucket + "/"
		if !strings.HasPrefix(withoutScheme, prefix) {
			return "", fmt.Errorf("object url bucket mismatch")
		}
		return strings.TrimPrefix(withoutScheme, prefix), nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(parsed.Path, "/"), nil
}
