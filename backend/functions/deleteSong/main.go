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
	var e error
	var song *repoTypes.Song
	var validate = utils.NewValidator()
	var songGUID = req.PathParameters["guid"]

	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if song, e = awsUtils.GetOneSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}
	if song.Location != "" {
		var songKey = utils.GetBucketKeyFromURL(song.Location)
		if e := awsUtils.DeleteObject(ctx, songKey, awsClients.S3); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if song.Image != "" {
		var imageKey = utils.GetBucketKeyFromURL(song.Image)
		if e := awsUtils.DeleteObject(ctx, imageKey, awsClients.S3); e != nil {
			return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
		}
	}
	if e := awsUtils.DeleteSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusNoContent, nil), nil
}

func main() {
	runtime.Start(handler)
}