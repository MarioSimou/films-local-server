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
	var validate = utils.NewValidator()
	var currentSong *repoTypes.Song
	var e error
	var m *utils.Multipart

	if e := utils.IsMultipartFormData(event.Headers["Content-Type"]); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if m, e = utils.ParseBodyToMultipart(event); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentSong, e = utils.GetOneSong(ctx, songGUID, dynamoDBClient); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var uploadOutput *manager.UploadOutput
	var songBucketKey = utils.GetBucketKey(repoTypes.SongType, currentSong.GUID, m.Ext)
	if uploadOutput, e = utils.UploadObject(ctx, songBucketKey, m, types.StorageClassOnezoneIa, s3Uploader); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.Location = uploadOutput.Location

	var newSong *repoTypes.Song
	if newSong, e = utils.PutSong(ctx, *currentSong, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
