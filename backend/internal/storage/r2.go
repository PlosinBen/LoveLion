// Package storage provides a thin wrapper around Cloudflare R2 (S3-compatible)
// so that handlers, workers and other services share one client configuration
// without depending on the S3 SDK types directly.
package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"lovelion/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage wraps an *s3.Client bound to a specific bucket + public domain.
type R2Storage struct {
	client       *s3.Client
	bucket       string
	publicDomain string
}

// NewR2Storage constructs an R2 client from the app config.
func NewR2Storage(cfg *config.Config) (*R2Storage, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithEndpointResolverWithOptions(resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKey, cfg.R2SecretKey, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("load R2 config: %w", err)
	}

	return &R2Storage{
		client:       s3.NewFromConfig(awsCfg),
		bucket:       cfg.R2Bucket,
		publicDomain: strings.TrimRight(cfg.R2PublicDomain, "/"),
	}, nil
}

// Client exposes the underlying S3 client for handlers that still need it.
func (r *R2Storage) Client() *s3.Client { return r.client }

// Bucket returns the configured bucket name.
func (r *R2Storage) Bucket() string { return r.bucket }

// PublicDomain returns the public CDN domain (without trailing slash).
func (r *R2Storage) PublicDomain() string { return r.publicDomain }

// PublicURL returns the fully-qualified URL for a key stored in the bucket.
func (r *R2Storage) PublicURL(key string) string {
	return r.publicDomain + "/" + strings.TrimLeft(key, "/")
}

// KeyFromURL strips the public domain prefix from a stored file URL.
// If the URL doesn't match the configured domain it is returned unchanged.
func (r *R2Storage) KeyFromURL(fullURL string) string {
	return strings.TrimPrefix(fullURL, r.publicDomain+"/")
}

// Upload stores an object at the given key. ContentType may be empty.
// Returns the full public URL.
func (r *R2Storage) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if _, err := r.client.PutObject(ctx, input); err != nil {
		return "", fmt.Errorf("r2 put %s: %w", key, err)
	}
	return r.PublicURL(key), nil
}

// Delete removes an object from the bucket. Missing objects are not an error
// from R2's perspective.
func (r *R2Storage) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2 delete %s: %w", key, err)
	}
	return nil
}

// Download fetches an object by key and returns its bytes and content-type.
func (r *R2Storage) Download(ctx context.Context, key string) ([]byte, string, error) {
	out, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", fmt.Errorf("r2 get %s: %w", key, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, "", fmt.Errorf("r2 read %s: %w", key, err)
	}

	var contentType string
	if out.ContentType != nil {
		contentType = *out.ContentType
	}
	return data, contentType, nil
}

// DownloadByURL is a convenience wrapper that extracts the key from a stored
// file URL and calls Download.
func (r *R2Storage) DownloadByURL(ctx context.Context, fullURL string) ([]byte, string, error) {
	return r.Download(ctx, r.KeyFromURL(fullURL))
}
