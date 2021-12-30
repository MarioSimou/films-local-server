package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

var (
	cfg            aws.Config
	dynamoDBClient *dynamodb.Client
)

type PostFilm struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

func init() {
	var e error

	if dynamoDBClient, e = utils.NewDynamoDBClient(context.Background()); e != nil {
		log.Fatalf("error: %v", e)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody PostFilm
	var now = time.Now()
	var validate = utils.NewValidator()
	var dynamoDBItem map[string]types.AttributeValue
	var e error

	if e := utils.DecodeEventBody(request.Body, &reqBody); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Struct(reqBody); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	var film = models.Film{
		GUID:        uuid.New().String(),
		Name:        reqBody.Name,
		Description: reqBody.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if dynamoDBItem, e = attributevalue.MarshalMap(film); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var putItemInput = dynamodb.PutItemInput{
		Item:      dynamoDBItem,
		TableName: aws.String("films"),
	}

	if _, e = dynamoDBClient.PutItem(ctx, &putItemInput); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, film), nil
}

func main() {
	runtime.Start(handler)
}
