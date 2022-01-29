package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/others/songs-local-server-sls/internal/packages/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrSongNotFound  = fmt.Errorf("err: song not found")
	ErrSongsNotFound = fmt.Errorf("err: songs not found")
)

func (c *AWSClients) GetSongByGUID(ctx context.Context, guid string) (*models.Song, error) {
	var song *models.Song
	var e error
	var guidMapMarshalled map[string]types.AttributeValue
	var guidMap = map[string]string{
		"GUID": guid,
	}

	if guidMapMarshalled, e = attributevalue.MarshalMap(guidMap); e != nil {
		return nil, e
	}

	var getItemInput = &dynamodb.GetItemInput{
		TableName: aws.String(models.SongsTableName),
		Key:       guidMapMarshalled,
	}
	var getItemOutput *dynamodb.GetItemOutput
	if getItemOutput, e = c.Dynamodb.GetItem(ctx, getItemInput); e != nil {
		return nil, e
	}

	if e := attributevalue.UnmarshalMap(getItemOutput.Item, &song); e != nil {
		return nil, e
	}

	if len(song.GUID) == 0 {
		return nil, ErrSongNotFound
	}

	return song, nil
}

func (c *AWSClients) PutSongByGUID(ctx context.Context, guid string, body interface{}) (*models.Song, error) {
	var e error
	var song *models.Song
	var bodyMarshalled map[string]types.AttributeValue

	if bodyMarshalled, e = attributevalue.MarshalMap(&body); e != nil {
		return nil, e
	}

	var putItemInput = &dynamodb.PutItemInput{
		TableName: aws.String(models.SongsTableName),
		Item:      bodyMarshalled,
	}

	if _, e := c.Dynamodb.PutItem(ctx, putItemInput); e != nil {
		return nil, e
	}

	if song, e = c.GetSongByGUID(ctx, guid); e != nil {
		return nil, e
	}

	return song, nil
}

func (c *AWSClients) DeleteSongByGUID(ctx context.Context, guid string) error {
	var e error
	var marshalledGUID types.AttributeValue

	if marshalledGUID, e = attributevalue.Marshal(guid); e != nil {
		return e
	}

	var putItemInput = &dynamodb.DeleteItemInput{
		TableName: aws.String(models.SongsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": marshalledGUID,
		},
	}

	if _, e := c.Dynamodb.DeleteItem(ctx, putItemInput); e != nil {
		return e
	}

	return nil
}

type GetSongsInput struct {
	Limit int32
}

func NewGetSongsInput(limit *int32) *GetSongsInput {
	if limit == nil {
		return &GetSongsInput{
			Limit: int32(20),
		}
	}

	return &GetSongsInput{
		Limit: *limit,
	}
}

func (c *AWSClients) GetSongs(ctx context.Context, input *GetSongsInput) (*[]models.Song, error) {
	var e error
	var scanOutput *dynamodb.ScanOutput
	var songs []models.Song
	var scanInput = &dynamodb.ScanInput{
		TableName: aws.String(models.SongsTableName),
		Limit:     &input.Limit,
	}

	if scanOutput, e = c.Dynamodb.Scan(ctx, scanInput); e != nil {
		return nil, e
	}

	if e := attributevalue.UnmarshalListOfMaps(scanOutput.Items, &songs); e != nil {
		return nil, e
	}

	if len(songs) == 0 {
		return nil, ErrSongsNotFound
	}

	return &songs, nil
}
