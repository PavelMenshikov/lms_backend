package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectStorage interface {
	GetPublicURL(ctx context.Context, key string) (string, error)
	UploadFile(ctx context.Context, file io.Reader, key string, size int64, mimeType string) (string, error)
}

type S3Client struct {
	S3          *s3.Client
	BucketName  string
	EndpointURL string
}

func NewS3Client(endpoint, region, bucket, accessKey, secretKey string) (*S3Client, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
				Source:        aws.EndpointSourceCustom,
			}, nil
		})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return &S3Client{
		S3: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
		BucketName:  bucket,
		EndpointURL: endpoint,
	}, nil
}

var _ ObjectStorage = (*S3Client)(nil)

func (c *S3Client) GetPublicURL(ctx context.Context, key string) (string, error) {
	return "/api/files/" + key, nil
}

func (c *S3Client) ServeFile(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/files/")
	if key == "" {
		http.Error(w, "file key is required", http.StatusBadRequest)
		return
	}

	result, err := c.S3.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &c.BucketName,
		Key:    &key,
	})
	if err != nil {
		slog.Error("S3 GetObject failed",
			slog.String("bucket", c.BucketName),
			slog.String("key", key),
			slog.String("error", err.Error()),
		)
		if strings.Contains(err.Error(), "NoSuchKey") {
			http.Error(w, "file not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	defer result.Body.Close()

	if result.ContentType != nil {
		w.Header().Set("Content-Type", *result.ContentType)
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	io.Copy(w, result.Body)
}

func (c *S3Client) UploadFile(ctx context.Context, file io.Reader, key string, size int64, mimeType string) (string, error) {
	_, err := c.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &c.BucketName,
		Key:           &key,
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   &mimeType,
	})

	if err != nil {
		slog.Error("S3 UploadFile failed",
			slog.String("bucket", c.BucketName),
			slog.String("key", key),
			slog.String("endpoint", c.EndpointURL),
			slog.Int64("size", size),
			slog.String("mime_type", mimeType),
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("failed to upload file to S3 (bucket=%s, key=%s, endpoint=%s): %w", c.BucketName, key, c.EndpointURL, err)
	}
	return key, nil
}
