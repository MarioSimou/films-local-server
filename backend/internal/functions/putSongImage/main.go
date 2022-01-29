package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
	"github.com/others/songs-local-server-sls/internal/packages/common"
	"github.com/others/songs-local-server-sls/internal/packages/models"
)

type Body struct {
	Image string `json:"image"`
}

func Handler(awsClients awsUtils.IAWSClients) common.LambdaHandler {
	return func(event common.Event) (common.Response, error) {
		fmt.Println("isbase64: ", event.IsBase64Encoded)
		var guid *string
		var e error
		var currentSong *models.Song
		var ctx, _ = common.NewContext(nil)
		var multipartResponse *common.MultipartResponse

		fmt.Println("SAMPLE: ", event.Body[:500])

		if guid, e = common.GetGUIDFromParameters(event.PathParameters); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		// var body Body
		// if e := common.ParseBody(event.Body, &body); e != nil {
		// 	return common.NewHTTPResponse(http.StatusBadRequest, e)
		// }
		// var bf, _ = base64.RawStdEncoding.DecodeString(body.Image)

		if multipartResponse, e = common.ParseMultipartContentType(event.Headers, event.Body); e != nil {
			return common.NewHTTPResponse(http.StatusBadRequest, e)
		}

		if currentSong, e = awsClients.GetSongByGUID(ctx, *guid); e != nil {
			if e == awsUtils.ErrSongNotFound {
				return common.NewHTTPResponse(http.StatusInternalServerError, e)
			}
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		var key = fmt.Sprintf("%s/image%s", currentSong.GUID, ".jpg")
		var href *string
		if href, e = awsClients.UploadOneObject(ctx, key, multipartResponse.Body); e != nil {
			return common.NewHTTPResponse(http.StatusInternalServerError, e)
		}

		currentSong.Href = *href

		var newSong *models.Song
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
