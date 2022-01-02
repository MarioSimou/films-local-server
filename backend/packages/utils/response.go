package utils

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func newResponse(status int, data interface{}) APIResponse {
	switch status {
	case http.StatusOK:
		return APIResponse{Status: status, Success: true, Data: data}
	case http.StatusNoContent:
		return APIResponse{Status: status, Success: true}
	default:
		if e, ok := data.(error); ok {
			return APIResponse{Status: status, Success: false, Message: e.Error()}
		}
		return APIResponse{Status: status, Success: false, Message: data.(string)}
	}
}

func NewAPIResponse(status int, data interface{}) events.APIGatewayProxyResponse {
	var body, _ = json.Marshal(newResponse(status, data))
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}
}
