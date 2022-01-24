package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/songs-local-server/backend/packages/awsUtils"
	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	awsClients *awsUtils.AWSClients
)

func init() {
	var e error
	var ctx = context.Background()

	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)

	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var scanInput = &dynamodb.ScanInput{TableName: aws.String(repoTypes.SongsTableName)}
	var scanOutput = &dynamodb.ScanOutput{}
	var e error

	if scanOutput, e = awsClients.DynamoDB.Scan(ctx, scanInput); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var songs []repoTypes.Song
	if e := attributevalue.UnmarshalListOfMaps(scanOutput.Items, &songs); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if len(songs) == 0 {
		return utils.NewAPIResponse(http.StatusNotFound, repoTypes.ErrSongNotFound), nil
	}

	return utils.NewAPIResponse(http.StatusOK, songs), nil
}

func main() {
	runtime.Start(handler)
}
