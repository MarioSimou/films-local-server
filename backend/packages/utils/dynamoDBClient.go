package utils

import (
	"context"
	"os"

	"github.com/MarioSimou/songs-local-server/backend/packages/models"
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

func GetOneSong(ctx context.Context, guid string, dynamoClient *dynamodb.Client) (*models.Song, error) {
	var guidKey, _ = attributevalue.Marshal(guid)
	var getSongOutput = &dynamodb.GetItemOutput{}
	var e error
	var song models.Song
	var getSongInput = &dynamodb.GetItemInput{
		TableName: aws.String(models.SongsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": guidKey,
		},
	}

	if getSongOutput, e = dynamoClient.GetItem(ctx, getSongInput); e != nil {
		return nil, e
	}

	if e := attributevalue.UnmarshalMap(getSongOutput.Item, &song); e != nil {
		return nil, e
	}

	if song.GUID == "" {
		return nil, ErrSongNotFound
	}

	return &song, nil
}

func PutSong(ctx context.Context, payload models.Song, dynamoDBClient *dynamodb.Client) (*models.Song, error) {
	var e error
	var item, _ = attributevalue.MarshalMap(payload)
	var putSongInput = &dynamodb.PutItemInput{
		TableName:    &models.SongsTableName,
		Item:         item,
		ReturnValues: types.ReturnValueNone,
	}

	if _, e = dynamoDBClient.PutItem(ctx, putSongInput); e != nil {
		return nil, e
	}

	var song *models.Song
	if song, e = GetOneSong(ctx, payload.GUID, dynamoDBClient); e != nil {
		return nil, e
	}

	return song, e
}

func DeleteSong(ctx context.Context, songGUID string, dynamoDBClient *dynamodb.Client) error {
	var songKey, _ = attributevalue.Marshal(songGUID)
	var deleteItemInput = dynamodb.DeleteItemInput{
		TableName: aws.String(models.SongsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": songKey,
		},
	}

	if _, e := dynamoDBClient.DeleteItem(ctx, &deleteItemInput); e != nil {
		return e
	}

	return nil
}
