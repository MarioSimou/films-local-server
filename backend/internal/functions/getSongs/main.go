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
		var songs *[]models.Song
		var ctx, _ = common.NewContext(nil)
		var limit *int32

		if limit, e = common.GetLimitFromQueryParameters(event.QueryStringParameters); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if songs, e = awsClients.GetSongs(ctx, awsUtils.NewGetSongsInput(limit)); e != nil {
			if e == awsUtils.ErrSongsNotFound {
				return common.NewHTTPResponse(http.StatusNotFound, e)
			}
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		return common.NewHTTPResponse(http.StatusOK, songs)
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
