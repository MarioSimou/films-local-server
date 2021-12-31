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
	s3Uploader     *manager.Uploader
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

	if s3Uploader, e = utils.NewS3Uploader(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var guid = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var currentFilm *models.Film
	var e error

	if e := utils.IsMultipartFormData(req.Headers["Content-Type"]); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Var(guid, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Var(req.Body, "required,multibyte"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentFilm, e = utils.GetOneFilm(ctx, guid, dynamoDBClient); e != nil {
		if e == utils.ErrFilmNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var subtitleKey = utils.GetSubtitleBucketKey(currentFilm.GUID)
	if currentFilm.Subtitle != "" {
		if e := utils.DeleteObject(ctx, subtitleKey, s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}

	var uploadOutput *manager.UploadOutput
	if uploadOutput, e = utils.UploadObject(ctx, subtitleKey, req.Body, types.StorageClassGlacierIr, s3Uploader); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentFilm.Subtitle = uploadOutput.Location

	var newFilm *models.Film
	if newFilm, e = utils.PutFilm(ctx, *currentFilm, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newFilm), nil
}

func main() {
	runtime.Start(handler)
}
