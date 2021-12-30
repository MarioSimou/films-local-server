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
	// 2) Unmarshal payload
	// 3) If the movie exists, update it, otherwise return 404

	return utils.NewAPIResponse(http.StatusOK, "update film"), nil
}

func main() {
	runtime.Start(handler)
}
