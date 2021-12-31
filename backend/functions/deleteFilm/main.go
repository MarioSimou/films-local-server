package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	dynamoDBClient *dynamodb.Client
	s3Client       *s3.Client
)

func init() {
	var e error
	var ctx = context.Background()
	if dynamoDBClient, e = utils.NewDynamoDBClient(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	if s3Client, e = utils.NewS3Client(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var e error
	var film *models.Film
	var validate = utils.NewValidator()
	var filmGUID = req.PathParameters["guid"]

	if e := validate.Var(filmGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if film, e = utils.GetOneFilm(ctx, filmGUID, dynamoDBClient); e != nil {
		if e == utils.ErrFilmNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var guidKey, _ = attributevalue.Marshal(film.GUID)
	var deleteItemInput = dynamodb.DeleteItemInput{
		TableName: aws.String(models.FilmsTableName),
		Key: map[string]types.AttributeValue{
			"GUID": guidKey,
		},
	}

	if _, e := dynamoDBClient.DeleteItem(ctx, &deleteItemInput); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	// check if we can merge them to one
	if film.Location != "" {
		if e := utils.DeleteObject(ctx, utils.GetFilmBucketKey(filmGUID), s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if film.Subtitle != "" {
		if e := utils.DeleteObject(ctx, utils.GetSubtitleBucketKey(filmGUID), s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if film.Image != "" {
		if e := utils.DeleteObject(ctx, utils.GetImageBucketKey(filmGUID, "jpg"), s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}

	return utils.NewAPIResponse(http.StatusNoContent, nil), nil
}

func main() {
	runtime.Start(handler)
}
