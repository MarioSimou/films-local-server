package main

import (
	"context"
	"log"
	"net/http"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
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
	var songGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var currentSong *repoTypes.Song
	var newSong *repoTypes.Song
	var e error

	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentSong, e = utils.GetOneSong(ctx, songGUID, dynamoDBClient); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if currentSong.Image == "" {
		return utils.NewAPIResponse(http.StatusNotFound, repoTypes.ErrImageNotFound), nil
	}

	var imageKey = utils.GetBucketKeyFromURL(currentSong.Image)
	if e := utils.DeleteObject(ctx, imageKey, s3Client); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.Image = ""
	if newSong, e = utils.PutSong(ctx, *currentSong, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
