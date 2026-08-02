package adapter

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryStorage struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryStorage(cloudinaryURL string) (*CloudinaryStorage, error) {
	if cloudinaryURL == "" {
		return nil, fmt.Errorf("CLOUDINARY_URL is empty")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &CloudinaryStorage{cld: cld}, nil
}

func (s *CloudinaryStorage) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "twitbet_avatars",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func (s *CloudinaryStorage) DeleteImage(ctx context.Context, imageURL string) error {
	if !strings.Contains(imageURL, "cloudinary.com") {
		return nil
	}

	uploadIdx := strings.Index(imageURL, "/upload/")
	if uploadIdx == -1 {
		return nil
	}

	path := imageURL[uploadIdx+len("/upload/"):]
	parts := strings.Split(path, "/")
	
	if len(parts) > 0 && strings.HasPrefix(parts[0], "v") {
		parts = parts[1:]
	}
	
	if len(parts) == 0 {
		return nil
	}

	publicIDWithExt := strings.Join(parts, "/")
	lastDot := strings.LastIndex(publicIDWithExt, ".")
	publicID := publicIDWithExt
	if lastDot != -1 {
		publicID = publicIDWithExt[:lastDot]
	}

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
