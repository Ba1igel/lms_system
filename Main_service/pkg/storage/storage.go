package storage

import (
	"context"
	"io"
)

type UploadInput struct {
	ObjectName  string
	Reader      io.Reader
	Size        int64
	ContentType string
}

type DownloadResult struct {
	Reader      io.ReadCloser
	Size        int64
	ContentType string
}

type FileStorage interface {
	Upload(ctx context.Context, input UploadInput) (string, error)
	Download(ctx context.Context, objectName string) (*DownloadResult, error)
}
