package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	var filmGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var currentFilm *models.Film
	var newFilm *models.Film
	var e error

	if e := validate.Var(filmGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentFilm, e = utils.GetOneFilm(ctx, filmGUID, dynamoDBClient); e != nil {
		if e == utils.ErrFilmNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if currentFilm.Subtitle == "" {
		return utils.NewAPIResponse(http.StatusNotFound, utils.ErrSubtitleNotFound), nil
	}

	var subtitleKey = utils.GetSubtitleBucketKey(currentFilm.GUID)

	if e := utils.DeleteObject(ctx, subtitleKey, s3Client); e != nil {
		// if e == utils.ErrS3ObjectNotFound {
		// 	return utils.NewAPIResponse(http.StatusNotFound, e), nil
		// }
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentFilm.Subtitle = ""
	if newFilm, e = utils.PutFilm(ctx, *currentFilm, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newFilm), nil
}

func main() {
	runtime.Start(handler)
}
