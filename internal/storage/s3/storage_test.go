package s3_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/greg-source/ukrainian-theater-events/internal/mock"
	"github.com/greg-source/ukrainian-theater-events/internal/storage/s3"
	"github.com/stretchr/testify/assert"
)

const testBucket = "test-bucket"

func TestStorage_Get(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		setup     func(c *mock.MockS3Client)
		want      []byte
		wantError bool
	}{
		{
			name: "success",
			path: "franko_theater/events.json",
			setup: func(c *mock.MockS3Client) {
				c.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(&s3sdk.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader(`{"key":"value"}`)),
				}, nil)
			},
			want: []byte(`{"key":"value"}`),
		},
		{
			name: "client error",
			path: "missing.json",
			setup: func(c *mock.MockS3Client) {
				c.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(nil, errors.New("no such key"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock.NewMockS3Client(ctrl)
			tt.setup(client)

			store := s3.New(client, testBucket)
			got, err := store.Get(context.Background(), tt.path)

			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStorage_Put(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		data      []byte
		setup     func(c *mock.MockS3Client)
		wantError bool
	}{
		{
			name: "success",
			path: "franko_theater/2026-05.json",
			data: []byte(`{"events":[]}`),
			setup: func(c *mock.MockS3Client) {
				c.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(&s3sdk.PutObjectOutput{}, nil)
			},
		},
		{
			name: "client error",
			path: "franko_theater/2026-05.json",
			data: []byte(`{"events":[]}`),
			setup: func(c *mock.MockS3Client) {
				c.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, errors.New("access denied"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock.NewMockS3Client(ctrl)
			tt.setup(client)

			store := s3.New(client, testBucket)
			err := store.Put(context.Background(), tt.path, tt.data)

			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
