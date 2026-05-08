package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Storage implements storage.Storage backed by an S3 bucket.
type Storage struct {
	client Client
	bucket string
}

// New creates a new S3-backed Storage for the given bucket.
func New(client Client, bucket string) *Storage {
	return &Storage{client: client, bucket: bucket}
}

// Get reads the contents of the object at the given path.
func (s *Storage) Get(ctx context.Context, path string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("getting object %s: %w", path, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("reading object %s: %w", path, err)
	}

	return data, nil
}

// Put writes data to the object at the given path.
func (s *Storage) Put(ctx context.Context, path string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("putting object %s: %w", path, err)
	}

	return nil
}
