package main

import (
	"context"
	"log"
	"net/http"
	"time"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var (
	cfg            aws.Config
	dynamoDBClient *dynamodb.Client
)

type PostSong struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

func init() {
	var e error
	var ctx = context.Background()

	if dynamoDBClient, e = utils.NewDynamoDBClient(ctx); e != nil {
		log.Fatalf("error: %v", e)
	}
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody PostSong
	var now = time.Now()
	var validate = utils.NewValidator()
	var newSong *repoTypes.Song
	var e error

	if e := utils.DecodeEventBody(event.Body, &reqBody, event.IsBase64Encoded); e != nil {
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

	if newSong, e = utils.PutSong(ctx, song, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
