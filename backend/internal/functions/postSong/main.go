package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
	"github.com/others/songs-local-server-sls/internal/packages/common"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

func Handler(awsClients awsUtils.IAWSClients) common.LambdaHandler {
	return func(event common.Event) (common.Response, error) {
		var e error
		var song *models.Song
		var body models.PostSong
		var ctx, _ = common.NewContext(nil)

		if e := common.ParseBody(event.Body, &body); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		body.GUID = uuid.NewString()
		body.CreatedAt = time.Now()

		if song, e = awsClients.PutSongByGUID(ctx, body.GUID, body); e != nil {
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
