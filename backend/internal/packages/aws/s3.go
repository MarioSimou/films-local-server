package aws

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

func (c *AWSClients) UploadOneObject(ctx context.Context, key string, body []byte) (*string, error) {
	var output *manager.UploadOutput
	var e error

	var input = &s3.PutObjectInput{
		Bucket:       &models.ServiceBucketName,
		ACL:          types.ObjectCannedACLPublicRead,
		Key:          &key,
		StorageClass: types.StorageClassStandard,
		Body:         bytes.NewReader(body),
	}

	if output, e = c.S3Uploader.Upload(ctx, input); e != nil {
		return nil, e
	}

	return &output.Location, nil
}

func (c *AWSClients) DeleteOneObject(ctx context.Context, key string) error {
	var input = &s3.DeleteObjectInput{
		Bucket: &models.ServiceBucketName,
		Key:    &key,
	}

	if _, e := c.S3.DeleteObject(ctx, input); e != nil {
		return e
	}
	return nil
}
