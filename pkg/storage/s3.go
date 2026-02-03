package storage

import (
	"context"
	"fmt"
	"io"

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
	return fmt.Sprintf("%s/%s/%s", c.EndpointURL, c.BucketName, key), nil
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
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return key, nil
}
