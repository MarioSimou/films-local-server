package utils

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	var env = os.Getenv("ENVIRONMENT")
	var cfg aws.Config
	var e error
	if env == "prod" {
		if cfg, e = config.LoadDefaultConfig(ctx); e != nil {
			return nil, e
		}
		return dynamodb.NewFromConfig(cfg), nil
	}

	var fn = func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: "http://dynamodb:4566"}, nil
	}

	if cfg, e = config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(fn)),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
				Source:          "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	); e != nil {
		return nil, e
	}

	return dynamodb.NewFromConfig(cfg), nil
}
