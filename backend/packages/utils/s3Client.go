package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	S3_BUCKET_NAME = "songs-backend-bucket"
)

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	var e error
	var cfg *aws.Config
	var s3Endpoint = os.Getenv("AWS_S3_ENDPOINT")
	var env repoTypes.EnvironmentType = os.Getenv("ENVIRONMENT")

	if env == repoTypes.Dev && len(s3Endpoint) == 0 {
		return nil, fmt.Errorf("error: environment '%s' with endpoint '%s'", env, s3Endpoint)
	}

	if cfg, e = GetAWSConfig(ctx, env, s3Endpoint); e != nil {
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
				return repoTypes.ErrS3ObjectNotFound
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
		StorageClass: storageClass,
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
