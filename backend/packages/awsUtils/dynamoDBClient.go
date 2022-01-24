package awsUtils

import (
	"context"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetOneSong(ctx context.Context, guid string, dynamoClient *dynamodb.Client) (*repoTypes.Song, error) {
	var guidKey, _ = attributevalue.Marshal(guid)
	var getSongOutput = &dynamodb.GetItemOutput{}
	var e error
	var song repoTypes.Song
	var getSongInput = &dynamodb.GetItemInput{
		TableName: aws.String(repoTypes.SongsTableName),
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
		return nil, repoTypes.ErrSongNotFound
	}

	return &song, nil
}

func PutSong(ctx context.Context, payload repoTypes.Song, dynamoDBClient *dynamodb.Client) (*repoTypes.Song, error) {
	var e error
	var item, _ = attributevalue.MarshalMap(payload)
	var putSongInput = &dynamodb.PutItemInput{
		TableName:    &repoTypes.SongsTableName,
		Item:         item,
		ReturnValues: types.ReturnValueNone,
	}

	if _, e = dynamoDBClient.PutItem(ctx, putSongInput); e != nil {
		return nil, e
	}

	var song *repoTypes.Song
	if song, e = GetOneSong(ctx, payload.GUID, dynamoDBClient); e != nil {
		return nil, e
	}

	return song, e
}

func DeleteSong(ctx context.Context, songGUID string, dynamoDBClient *dynamodb.Client) error {
	var songKey, _ = attributevalue.Marshal(songGUID)
	var deleteItemInput = dynamodb.DeleteItemInput{
		TableName: aws.String(repoTypes.SongsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": songKey,
		},
	}

	if _, e := dynamoDBClient.DeleteItem(ctx, &deleteItemInput); e != nil {
		return e
	}

	return nil
}
