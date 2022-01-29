package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
	"github.com/others/songs-local-server-sls/internal/packages/common"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

func Handler(awsClients awsUtils.IAWSClients) common.LambdaHandler {
	return func(event common.Event) (common.Response, error) {
		var e error
		var song *models.Song
		var ctx, _ = common.NewContext(nil)
		var guid *string

		if guid, e = common.GetGUIDFromParameters(event.PathParameters); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if song, e = awsClients.GetSongByGUID(ctx, *guid); e != nil {
			if e == awsUtils.ErrSongNotFound {
				return common.NewHTTPResponse(http.StatusNotFound, e)
			}
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		return common.NewHTTPResponse(http.StatusOK, song)
	}
}

func main() {
	var ctx = context.Background()
	var e error
	var awsClients awsUtils.IAWSClients
	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("error: %v\n", e)
	}

	lambda.Start(common.WithAWS(ctx, awsClients, Handler))
}
