package storage

import (
	"context"
	"io"
	"time"
)

type UploadOptions struct {
	Public bool
}

type StoredObject struct {
	Key       string `json:"key"`
	PublicURL string `json:"public_url,omitempty"`
}

type Adapter interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, opts UploadOptions) (StoredObject, error)
	Delete(ctx context.Context, key string) error
	GenerateSignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	PublicURL(key string) string
}
