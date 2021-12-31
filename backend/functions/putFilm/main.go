package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MarioSimou/films-local-server/backend/packages/models"
	"github.com/MarioSimou/films-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jinzhu/copier"
)

var (
	dynamoDBClient *dynamodb.Client
)

type bodyBinding struct {
	Name        string    `json:"name" validate:"required_without_all=Description"`
	Description string    `json:"description" validate:"max=1000"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func init() {
	var e error
	if dynamoDBClient, e = utils.NewDynamoDBClient(context.Background()); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var filmGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var body bodyBinding

	if e := utils.DecodeEventBody(req.Body, &body); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Var(filmGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Struct(body); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	body.UpdatedAt = time.Now()

	var currentFilm *models.Film
	var e error
	if currentFilm, e = utils.GetOneFilm(ctx, filmGUID, dynamoDBClient); e != nil {
		if e == utils.ErrFilmNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if e := copier.CopyWithOption(currentFilm, &body, copier.Option{IgnoreEmpty: true}); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var item, _ = attributevalue.MarshalMap(*currentFilm)
	var putItemInput = dynamodb.PutItemInput{
		TableName: &models.FilmsTableName,
		Item:      item,
	}

	if _, e := dynamoDBClient.PutItem(ctx, &putItemInput); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	var newFilm *models.Film
	if newFilm, e = utils.GetOneFilm(ctx, filmGUID, dynamoDBClient); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newFilm), nil
}

func main() {
	runtime.Start(handler)
}
