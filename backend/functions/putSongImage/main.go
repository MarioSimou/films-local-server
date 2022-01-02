package main

import (
	"context"
	"log"
	"net/http"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var songGUID = event.PathParameters["guid"]
	var currentSong *repoTypes.Song
	var e error
	var m *utils.Multipart

	if e := utils.IsMultipartFormData(event.Headers["Content-Type"]); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if m, e = utils.ParseBodyToMultipart(event); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentSong, e = utils.GetOneSong(ctx, songGUID, dynamoDBClient); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var imageBucketKey = utils.GetBucketKey(repoTypes.ImageType, currentSong.GUID, m.Ext)
	if currentSong.Image != "" {
		if e := utils.DeleteObject(ctx, imageBucketKey, s3Client); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}

	var uploadOutput *manager.UploadOutput
	if uploadOutput, e = utils.UploadObject(ctx, imageBucketKey, m, types.StorageClassStandard, s3Uploader); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.Image = uploadOutput.Location

	var newSong *repoTypes.Song
	if newSong, e = utils.PutSong(ctx, *currentSong, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
