package storage

import (
	"fmt"
	"time"

	"github.com/soundmarket/backend/internal/config"
)

type Adapter interface {
	SignedURL(objectKey string) string
}

type S3Adapter struct {
	cfg *config.Config
}

func NewS3Adapter(cfg *config.Config) *S3Adapter {
	return &S3Adapter{cfg: cfg}
}

func (a *S3Adapter) SignedURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s?expires=%d", a.cfg.S3Endpoint, a.cfg.S3Bucket, objectKey, time.Now().Add(a.cfg.SignedURLTTL).Unix())
}

