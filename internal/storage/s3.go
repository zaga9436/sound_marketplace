package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/soundmarket/backend/internal/config"
)

type S3Adapter struct {
	cfg       *config.Config
	client    *s3.Client
	presigner *s3.PresignClient
}

func NewS3Adapter(cfg *config.Config) (*S3Adapter, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:               strings.TrimRight(normalizeEndpoint(cfg.S3Endpoint, cfg.S3UseSSL), "/"),
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.S3ForcePathStyle
	})

	return &S3Adapter{
		cfg:       cfg,
		client:    client,
		presigner: s3.NewPresignClient(client),
	}, nil
}

func (a *S3Adapter) Upload(ctx context.Context, key, contentType string, body io.Reader, opts UploadOptions) (StoredObject, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(a.cfg.S3Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}
	if opts.Public {
		input.ACL = s3types.ObjectCannedACLPublicRead
	}

	if _, err := a.client.PutObject(ctx, input); err != nil {
		return StoredObject{}, err
	}

	object := StoredObject{Key: key}
	if opts.Public {
		object.PublicURL = a.PublicURL(key)
	}
	return object, nil
}

func (a *S3Adapter) Delete(ctx context.Context, key string) error {
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.cfg.S3Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (a *S3Adapter) GenerateSignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	resp, err := a.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.cfg.S3Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = ttl
	})
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

func (a *S3Adapter) PublicURL(key string) string {
	base := strings.TrimRight(normalizeEndpoint(a.cfg.S3Endpoint, a.cfg.S3UseSSL), "/")
	if a.cfg.S3ForcePathStyle {
		return fmt.Sprintf("%s/%s/%s", base, a.cfg.S3Bucket, strings.TrimLeft(key, "/"))
	}

	parsed, err := url.Parse(base)
	if err != nil || parsed.Host == "" {
		return fmt.Sprintf("%s/%s", base, strings.TrimLeft(key, "/"))
	}
	return fmt.Sprintf("%s://%s.%s/%s", parsed.Scheme, a.cfg.S3Bucket, parsed.Host, strings.TrimLeft(key, "/"))
}

func BuildObjectKey(prefix, filename string) string {
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	return path.Join(prefix, filename)
}

func normalizeEndpoint(endpoint string, useSSL bool) string {
	endpoint = strings.TrimSpace(endpoint)
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}
	scheme := "http://"
	if useSSL {
		scheme = "https://"
	}
	return scheme + endpoint
}
