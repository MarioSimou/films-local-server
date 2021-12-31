package utils

import (
	"context"
	"os"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetAWSConfig(ctx context.Context, endpoint string) (*aws.Config, error) {
	var env = os.Getenv("ENVIRONMENT")
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

func NewDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	var e error
	var cfg *aws.Config

	if cfg, e = GetAWSConfig(ctx, "http://dynamodb:4566"); e != nil {
		return nil, e
	}

	return dynamodb.NewFromConfig(*cfg), nil
}

func GetOneFilm(ctx context.Context, guid string, dynamoClient *dynamodb.Client) (*models.Film, error) {
	var guidKey, _ = attributevalue.Marshal(guid)
	var getFilmInput = &dynamodb.GetItemInput{
		TableName: aws.String(models.FilmsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": guidKey,
		},
	}
	var getFilmOutput = &dynamodb.GetItemOutput{}
	var e error
	var film models.Film

	if getFilmOutput, e = dynamoClient.GetItem(ctx, getFilmInput); e != nil {
		return nil, e
	}

	if e := attributevalue.UnmarshalMap(getFilmOutput.Item, &film); e != nil {
		return nil, e
	}

	if film.GUID == "" {
		return nil, ErrFilmNotFound
	}

	return &film, nil
}

func PutFilm(ctx context.Context, payload models.Film, dynamoDBClient *dynamodb.Client) (*models.Film, error) {
	var e error
	var item, _ = attributevalue.MarshalMap(payload)
	var putItemInput = &dynamodb.PutItemInput{
		TableName:    &models.FilmsTableName,
		Item:         item,
		ReturnValues: types.ReturnValueNone,
	}

	if _, e = dynamoDBClient.PutItem(ctx, putItemInput); e != nil {
		return nil, e
	}

	var film *models.Film
	if film, e = GetOneFilm(ctx, payload.GUID, dynamoDBClient); e != nil {
		return nil, e
	}

	return film, e
}
