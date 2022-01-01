package utils

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	S3_BUCKET_NAME = "songs-local-server"
)

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	var e error
	var cfg *aws.Config

	if cfg, e = GetAWSConfig(ctx, "http://s3:4566"); e != nil {
		return nil, e
	}

	var fn = func(options *s3.Options) {
		options.UsePathStyle = true
	}
	return s3.NewFromConfig(*cfg, fn), nil
}

func NewS3Uploader(ctx context.Context, fns ...func(*manager.Uploader)) (*manager.Uploader, error) {
	var s3Client *s3.Client
	var e error

	if s3Client, e = NewS3Client(ctx); e != nil {
		return nil, e
	}

	return manager.NewUploader(s3Client, fns...), nil
}

func HeadObject(ctx context.Context, key string, s3Client *s3.Client) error {
	var e error
	var headObjectInput = &s3.HeadObjectInput{
		Key:    &key,
		Bucket: &S3_BUCKET_NAME,
	}

	if _, e = s3Client.HeadObject(ctx, headObjectInput); e != nil {
		var apiErr smithy.APIError
		if errors.As(e, &apiErr); e != nil {
			if apiErr.ErrorCode() == "NotFound" {
				return ErrS3ObjectNotFound
			}
		}
		return e
	}
	return nil
}

func DeleteObject(ctx context.Context, key string, s3Client *s3.Client) error {
	var e error
	var deleteObjectInput = &s3.DeleteObjectInput{
		Key:    &key,
		Bucket: &S3_BUCKET_NAME,
	}

	if e := HeadObject(ctx, key, s3Client); e != nil {
		return e
	}

	if _, e = s3Client.DeleteObject(ctx, deleteObjectInput); e != nil {
		return e
	}

	return nil
}

func UploadObject(ctx context.Context, key string, m *Multipart, storageClass types.StorageClass, uploader *manager.Uploader) (*manager.UploadOutput, error) {
	var putObjectInput = &s3.PutObjectInput{
		Bucket:       &S3_BUCKET_NAME,
		Key:          &key,
		Body:         bytes.NewReader(m.Body),
		StorageClass: types.StorageClassStandard,
	}

	return uploader.Upload(ctx, putObjectInput)
}

func PutObject(ctx context.Context, key string, body string, s3Client *s3.Client) error {
	var putObjectInput = &s3.PutObjectInput{
		Bucket: &S3_BUCKET_NAME,
		Key:    &key,
		Body:   strings.NewReader(body),
	}

	if _, e := s3Client.PutObject(ctx, putObjectInput); e != nil {
		return e
	}

	return nil
}
