package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/copier"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
	"github.com/others/songs-local-server-sls/internal/packages/common"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

func Handler(awsClients awsUtils.IAWSClients) common.LambdaHandler {
	return func(event common.Event) (common.Response, error) {
		var newSong *models.Song
		var e error
		var body models.PutSong
		var currentSong *models.Song
		var guid *string
		var ctx, _ = common.NewContext(nil)

		if guid, e = common.GetGUIDFromParameters(event.PathParameters); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if e := common.ParseBody(event.Body, &body); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if currentSong, e = awsClients.GetSongByGUID(ctx, *guid); e != nil {
			if e == awsUtils.ErrSongNotFound {
				return common.NewHTTPResponse(http.StatusNotFound, e)
			}
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if e := copier.CopyWithOption(currentSong, &body, copier.Option{IgnoreEmpty: true}); e != nil {
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		if newSong, e = awsClients.PutSongByGUID(ctx, currentSong.GUID, currentSong); e != nil {
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		return common.NewHTTPResponse(http.StatusOK, newSong)
	}
}

func main() {
	var ctx = context.Background()
	var awsClients awsUtils.IAWSClients
	var e error
	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("error: %v\n", e)
	}

	lambda.Start(common.WithAWS(ctx, awsClients, Handler))
}
