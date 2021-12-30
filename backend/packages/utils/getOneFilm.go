package utils

import (
	"context"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
	return &film, nil
}
