package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage implements S3-compatible storage for Cloudflare R2
type R2Storage struct {
	client     *s3.Client
	bucket     string
	publicURL  string
}

// R2Config holds the configuration for R2 storage
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string // Optional public URL for accessing files
	Region          string // Default: auto
}

// NewR2Storage creates a new R2 storage instance
func NewR2Storage(cfg *R2Config) (*R2Storage, error) {
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		return nil, fmt.Errorf("missing required R2 configuration")
	}

	// Set default region if not provided
	region := cfg.Region
	if region == "" {
		region = "auto"
	}

	// Create AWS configuration for R2
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom endpoint for R2
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	return &R2Storage{
		client:    client,
		bucket:    cfg.BucketName,
		publicURL: cfg.PublicURL,
	}, nil
}

// UploadFile uploads a file to R2
func (r *R2Storage) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	// Upload to R2
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to R2: %w", err)
	}

	// Return the public URL
	return r.GetPublicURL(key), nil
}

// UploadFileFromReader uploads a file from an io.Reader to R2
func (r *R2Storage) UploadFileFromReader(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	// Read all data from reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file data: %w", err)
	}

	return r.UploadFile(ctx, key, data, contentType)
}

// DeleteFile deletes a file from R2
func (r *R2Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %w", err)
	}

	return nil
}

// GetFile retrieves a file from R2
func (r *R2Storage) GetFile(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from R2: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

// FileExists checks if a file exists in R2
func (r *R2Storage) FileExists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetPublicURL returns the public URL for a file
func (r *R2Storage) GetPublicURL(key string) string {
	if r.publicURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(r.publicURL, "/"), key)
	}
	// Default R2 public URL format
	return fmt.Sprintf("https://pub-%s.r2.dev/%s", r.bucket, key)
}

// ListFiles lists all files in the bucket with a given prefix
func (r *R2Storage) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	result, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

// CopyFile copies a file within R2
func (r *R2Storage) CopyFile(ctx context.Context, sourceKey, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", r.bucket, sourceKey)
	_, err := r.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(r.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destKey),
	})
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
