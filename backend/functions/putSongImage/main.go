package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/songs-local-server/backend/packages/awsUtils"
	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var (
	awsClients *awsUtils.AWSClients
)

func init() {
	var e error
	var ctx = context.Background()

	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var songGUID = event.PathParameters["guid"]
	var currentSong *repoTypes.Song
	var e error
	var m *utils.Multipart
	var newSong *repoTypes.Song

	if e := utils.IsMultipartFormData(event.Headers["Content-Type"]); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if m, e = utils.ParseBodyToMultipart(event); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentSong, e = awsUtils.GetOneSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var imageBucketKey = utils.GetBucketKey(repoTypes.ImageType, currentSong.GUID, m.Ext)
	if _, e := awsUtils.UploadOne(ctx, awsClients.SNS, imageBucketKey, repoTypes.ImageField, m.Body); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.ImageStatus = repoTypes.Pending

	if newSong, e = awsUtils.PutSong(ctx, *currentSong, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
