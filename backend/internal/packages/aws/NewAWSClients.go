package aws

import (
	"context"
	"os"

	"github.com/others/songs-local-server-sls/internal/packages/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSClients struct {
	Dynamodb   *dynamodb.Client
	S3         *s3.Client
	S3Uploader *manager.Uploader
}

type IAWSClients interface {
	PutSongByGUID(ctx context.Context, guid string, body interface{}) (*models.Song, error)
	DeleteSongByGUID(ctx context.Context, guid string) error
	GetSongByGUID(ctx context.Context, guid string) (*models.Song, error)
	GetSongs(ctx context.Context, input *GetSongsInput) (*[]models.Song, error)
	UploadOneObject(ctx context.Context, key string, body []byte) (*string, error)
	DeleteOneObject(ctx context.Context, key string) error
}

func newConfig(ctx context.Context) (aws.Config, error) {
	var awsEndpoint = os.Getenv("AWS_ENDPOINT")

	var endpointResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if len(awsEndpoint) > 0 {
			return aws.Endpoint{
				URL: awsEndpoint,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	return config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(endpointResolver))
}

func NewAWSClients(ctx context.Context) (IAWSClients, error) {
	var cfg aws.Config
	var e error

	if cfg, e = newConfig(ctx); e != nil {
		return nil, e
	}

	var s3ClientFn = func(opts *s3.Options) {
		opts.UsePathStyle = true
	}
	var s3Client = s3.NewFromConfig(cfg, s3ClientFn)
	return &AWSClients{
		Dynamodb:   dynamodb.NewFromConfig(cfg),
		S3:         s3Client,
		S3Uploader: manager.NewUploader(s3Client),
	}, nil
}
