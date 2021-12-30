package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
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
	var e error
	var film *models.Film
	var validate = utils.NewValidator()
	var guid = req.PathParameters["guid"]

	if e := validate.Var(guid, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if film, e = utils.GetOneFilm(ctx, guid, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if film.GUID == "" {
		return utils.NewAPIResponse(http.StatusNotFound, utils.ErrDocNotFound), nil
	}

	return utils.NewAPIResponse(http.StatusOK, film), nil
}

func main() {
	runtime.Start(handler)
}
