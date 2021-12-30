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
	// 1) Find existing film
	// 2) If film doesn't exist, return a 404
	// 3) Delete subtitle of the film from s3 bucket
	// 4) Update dynamodb to capture the latest state

	return utils.NewAPIResponse(http.StatusOK, "delete subtitle"), nil
}

func main() {
	runtime.Start(handler)
}
