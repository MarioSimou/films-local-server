package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/songs-local-server/backend/packages/models"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
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
	var e error
	var song *models.Song
	var validate = utils.NewValidator()
	var songGUID = req.PathParameters["guid"]

	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if song, e = utils.GetOneSong(ctx, songGUID, dynamoDBClient); e != nil {
		if e == utils.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}
	if song.Location != "" {
		var songKey = utils.GetBucketKeyFromURL(song.Location)
		if e := utils.DeleteObject(ctx, songKey, s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if song.Image != "" {
		var imageKey = utils.GetBucketKeyFromURL(song.Image)
		if e := utils.DeleteObject(ctx, imageKey, s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if e := utils.DeleteSong(ctx, songGUID, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusNoContent, nil), nil
}

func main() {
	runtime.Start(handler)
}
