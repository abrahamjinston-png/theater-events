package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client defines the interface for S3 operations used by the storage.
//
//go:generate mockgen -source=client.go -destination=../../mock/mock_s3_client.go -package=mock -mock_names Client=MockS3Client
type Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// NewClient creates a new S3 client with the specified region.
func NewClient(ctx context.Context, region string) (Client, error) {
	if region == "" {
		return nil, fmt.Errorf("AWS region must be specified")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}
