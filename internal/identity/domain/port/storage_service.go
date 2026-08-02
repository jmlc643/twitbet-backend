package port

import (
	"context"
	"io"
)

type StorageService interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
	DeleteImage(ctx context.Context, imageURL string) error
}
