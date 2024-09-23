package services

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/ilholbea/viego-art-backend/types"
	"github.com/sirupsen/logrus"
	"time"
)

type ImageService interface {
	UploadImage(ctx context.Context, bucket string, file *types.UploadedFile) error
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (is *imageService) UploadImage(ctx context.Context, bucket string, file *types.UploadedFile) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	// Create a GCS writer
	wc := client.Bucket(bucket).Object(file.Filename).NewWriter(ctx)
	wc.ContentType = file.ContentType

	// Write file content to GCS
	if _, err := wc.Write([]byte(file.File)); err != nil {
		return fmt.Errorf("error writing to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("error closing GCS writer: %w", err)
	}

	logrus.Infof("Blob %v uploaded to bucket %v.", file.Filename, bucket)
	return nil
}
