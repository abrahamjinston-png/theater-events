package storage

import "context"

// Storage abstracts read/write access to a file-like store (S3, local filesystem, etc.).
//
//go:generate mockgen -source=storage.go -destination=../mock/mock_storage.go -package=mock
type Storage interface {
	// Get reads the contents at the given path.
	Get(ctx context.Context, path string) ([]byte, error)

	// Put writes data to the given path, creating or overwriting it.
	Put(ctx context.Context, path string, data []byte) error
}
