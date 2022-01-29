package common

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	awsUtils "github.com/others/songs-local-server-sls/internal/packages/aws"
)

type Response events.APIGatewayProxyResponse
type Event events.APIGatewayProxyRequest
type LambdaHandler func(event Event) (Response, error)

type ResponseBody struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WithAWS(cxt context.Context, awsClients awsUtils.IAWSClients, lambda func(awsClients awsUtils.IAWSClients) LambdaHandler) LambdaHandler {
	return lambda(awsClients)
}

func NewResponseBody(status int, data interface{}) ResponseBody {
	if e, ok := data.(error); ok {
		return ResponseBody{
			Status:  status,
			Success: false,
			Message: e.Error(),
		}
	}

	return ResponseBody{
		Status:  status,
		Success: true,
		Data:    data,
	}
}

func NewHTTPResponse(status int, data interface{}) (Response, error) {
	var bf []byte
	var e error
	var body = NewResponseBody(status, data)

	if bf, e = json.Marshal(body); e != nil {
		return Response{
			StatusCode: http.StatusBadRequest,
			Body:       e.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		}, nil
	}

	return Response{
		StatusCode: status,
		Body:       string(bf),
	}, nil
}
