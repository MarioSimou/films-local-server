package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MarioSimou/songs-local-server/backend/packages/awsUtils"
	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var (
	awsClients *awsUtils.AWSClients
)

type PostSong struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

func init() {
	var e error
	var ctx = context.Background()

	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)

	}
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody PostSong
	var now = time.Now()
	var validate = utils.NewValidator()
	var newSong *repoTypes.Song
	var e error

	if e := utils.DecodeEventBody(event.Body, &reqBody); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Struct(reqBody); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	var song = repoTypes.Song{
		GUID:        uuid.New().String(),
		Name:        reqBody.Name,
		Description: reqBody.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if newSong, e = awsUtils.PutSong(ctx, song, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
