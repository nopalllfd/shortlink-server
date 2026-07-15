package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectStorage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type R2Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewR2Storage(
	client *s3.Client,
	bucket string,
	publicURL string,
) *R2Storage {
	return &R2Storage{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}
}

func (r *R2Storage) Upload(
	ctx context.Context,
	key string,
	data []byte,
	contentType string,
) (string, error) {

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", r.publicURL, key), nil
}

func (r *R2Storage) Delete(
	ctx context.Context,
	key string,
) error {

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})

	return err
}
