package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	dynamoDBClient *dynamodb.Client
	s3Client       *s3.Client
	uploader       *manager.Uploader
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
	if uploader, e = utils.NewS3Uploader(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var filmGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var currentFilm *models.Film
	var e error

	if e := utils.IsMultipartFormData(req.Headers["Content-Type"]); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if e := validate.Var(filmGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if e := validate.Var(req.Body, "required,multibyte"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if currentFilm, e = utils.GetOneFilm(ctx, filmGUID, dynamoDBClient); e != nil {
		if e == utils.ErrFilmNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var filmBucketKey = utils.GetFilmBucketKey(currentFilm.GUID)
	if currentFilm.Location != "" {
		if e := utils.DeleteObject(ctx, filmBucketKey, s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}

	var uploadOutput *manager.UploadOutput
	if uploadOutput, e = utils.UploadObject(ctx, filmBucketKey, req.Body, types.StorageClassGlacierIr, uploader); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentFilm.Location = uploadOutput.Location
	var film *models.Film
	if film, e = utils.PutFilm(ctx, *currentFilm, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, film), nil
}

func main() {
	runtime.Start(handler)
}
