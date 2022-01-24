package awsUtils

import (
	"context"
	"fmt"
	"os"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type AWSClients struct {
	S3         *s3.Client
	DynamoDB   *dynamodb.Client
	SNS        *sns.Client
	S3Uploader *manager.Uploader
}

func GetAWSConfig(ctx context.Context, env repoTypes.EnvironmentType, endpoint string) (*aws.Config, error) {
	var cfg aws.Config
	var e error

	if env == "prod" {
		if cfg, e = config.LoadDefaultConfig(ctx); e != nil {
			return nil, e
		}
		return &cfg, e
	}

	var customResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			SigningRegion: "us-east-1",
			URL:           endpoint,
		}, nil
	})

	if cfg, e = config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(customResolver),
	); e != nil {
		return nil, e
	}

	return &cfg, nil
}

func NewAWSClients(ctx context.Context) (*AWSClients, error) {
	var e error
	var cfg *aws.Config
	var awsEndpoint = os.Getenv("AWS_ENDPOINT")
	var env repoTypes.EnvironmentType = os.Getenv("ENVIRONMENT")

	if env == repoTypes.Dev && len(awsEndpoint) == 0 {
		return nil, fmt.Errorf("error: environment '%s' with endpoint '%s'", env, awsEndpoint)
	}

	if cfg, e = GetAWSConfig(ctx, env, awsEndpoint); e != nil {
		return nil, e
	}

	var s3ConfigFn = func(options *s3.Options) {
		options.UsePathStyle = true
	}
	var s3Client = s3.NewFromConfig(*cfg, s3ConfigFn)
	var s3UploaderClient = manager.NewUploader(s3Client)
	var snsClient = sns.NewFromConfig(*cfg)
	var dynamodbClient = dynamodb.NewFromConfig(*cfg)

	return &AWSClients{
		S3:         s3Client,
		DynamoDB:   dynamodbClient,
		SNS:        snsClient,
		S3Uploader: s3UploaderClient,
	}, nil
}
