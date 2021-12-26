package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Status  int
	Success bool
	Message string
	Data    interface{}
}

func newResponse(status int, data interface{}) Response {
	switch status {
	case http.StatusOK:
		return Response{Status: status, Success: true, Data: data}
	case http.StatusNoContent:
		return Response{Status: status, Success: true}
	default:
		return Response{Status: status, Success: false, Message: data.(string)}
	}
}

func newAPIResponse(status int, data interface{}) events.APIGatewayProxyResponse {
	var body, _ = json.Marshal(newResponse(status, data))
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}
}

func handler() (events.APIGatewayProxyResponse, error) {
	return newAPIResponse(http.StatusOK, "post film"), nil
}

func main() {
	runtime.Start(handler)
}
