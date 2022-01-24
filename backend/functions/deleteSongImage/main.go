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

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var songGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var currentSong *repoTypes.Song
	var newSong *repoTypes.Song
	var e error

	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if currentSong, e = awsUtils.GetOneSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if currentSong.Image == "" {
		return utils.NewAPIResponse(http.StatusNotFound, repoTypes.ErrImageNotFound), nil
	}

	var imageKey = utils.GetBucketKeyFromURL(currentSong.Image)
	if e := awsUtils.DeleteObject(ctx, imageKey, awsClients.S3); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.Image = ""
	if newSong, e = awsUtils.PutSong(ctx, *currentSong, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
