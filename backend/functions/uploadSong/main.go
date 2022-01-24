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

	if currentSong, e = awsUtils.GetOneSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var songBucketKey = utils.GetBucketKey(repoTypes.SongType, currentSong.GUID, m.Ext)

	if _, e := awsUtils.UploadOne(ctx, awsClients.SNS, songBucketKey, repoTypes.LocatioField, m.Body); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	currentSong.LocationStatus = repoTypes.Pending

	var newSong *repoTypes.Song
	if newSong, e = awsUtils.PutSong(ctx, *currentSong, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
