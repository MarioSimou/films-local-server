package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dynamoDBClient *dynamodb.Client
)

func init() {
	var e error
	if dynamoDBClient, e = utils.NewDynamoDBClient(context.Background()); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1) Find existing movie
	// 2) If film doesn't exist in dynamodb, return a 404
	// 3) if the movie already has a film registered to it, delete it
	// 4) Upload new film of the movie
	// 5) Update dynamodb to include the latest state
	return utils.NewAPIResponse(http.StatusOK, "upload film"), nil
}

func main() {
	runtime.Start(handler)
}
