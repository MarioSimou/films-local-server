package awsUtils

import (
	"bytes"
	"context"
	"errors"
	"strings"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	S3_BUCKET_NAME = "songs-backend-bucket"
)

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

func UploadObject(ctx context.Context, key string, body []byte, storageClass types.StorageClass, uploader *manager.Uploader) (*manager.UploadOutput, error) {
	var putObjectInput = &s3.PutObjectInput{
		Bucket:       &S3_BUCKET_NAME,
		Key:          &key,
		Body:         bytes.NewReader(body),
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
