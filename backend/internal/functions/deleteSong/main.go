package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
	"github.com/others/songs-local-server-sls/internal/packages/common"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

func Handler(awsClients awsUtils.IAWSClients) common.LambdaHandler {
	return func(event common.Event) (common.Response, error) {
		var e error
		var currentSong *models.Song
		var ctx, _ = common.NewContext(nil)
		var guid *string

		if guid, e = common.GetGUIDFromParameters(event.PathParameters); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if currentSong, e = awsClients.GetSongByGUID(ctx, *guid); e != nil {
			if e == awsUtils.ErrSongNotFound {
				return common.NewHTTPResponse(http.StatusNotFound, e)
			}
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		if len(currentSong.Image) > 0 {
			var imageKey = fmt.Sprintf("%s/image%s", currentSong.GUID, filepath.Ext(currentSong.Image))
			if e := awsClients.DeleteOneObject(ctx, imageKey); e != nil {
				return common.NewHTTPResponse(http.StatusInternalServerError, nil)
			}
		}

		if len(currentSong.Href) > 0 {
			var songKey = fmt.Sprintf("%s/song%s", currentSong.GUID, filepath.Ext(currentSong.Href))
			if e := awsClients.DeleteOneObject(ctx, songKey); e != nil {
				return common.NewHTTPResponse(http.StatusInternalServerError, nil)
			}
		}

		if e := awsClients.DeleteSongByGUID(ctx, currentSong.GUID); e != nil {
			return common.NewHTTPResponse(http.StatusInternalServerError, nil)
		}

		return common.NewHTTPResponse(http.StatusNoContent, nil)
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
