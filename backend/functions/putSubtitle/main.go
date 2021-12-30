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
	// 3) If the film already has a subtitel, delete it from s3 bucket
	// 4) Upload the latest subtitle to the s3 bucket
	// 5) Update dynamodb to capture the latest state

	return utils.NewAPIResponse(http.StatusOK, "update subtitle of a film"), nil
}

func main() {
	runtime.Start(handler)
}
