package mocks

import (
	"context"
	"io"
	"lms_backend/pkg/storage"
)

type S3StorageMock struct{}

func NewS3StorageMock() *S3StorageMock {
	return &S3StorageMock{}
}

var _ storage.ObjectStorage = (*S3StorageMock)(nil)

func (m *S3StorageMock) GetPublicURL(ctx context.Context, key string) (string, error) {
	return "https://mock-s3-url.com/" + key, nil
}

func (m *S3StorageMock) UploadFile(ctx context.Context, file io.Reader, key string, size int64, mimeType string) (string, error) {
	return key, nil
}
