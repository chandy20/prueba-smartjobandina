package lib

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// EmptyResponse func for response
func EmptyResponse(statusCode int) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	return resp
}

//ResponseError function to build json error
func ResponseError(statusCode int, err error) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]interface{}{
		"message": err.Error(),
	})

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)
	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	return resp
}

// JSONResponse receives a JSON body and a code and returns a Response of type APIGatewayProxyResponse
func JSONResponse(statusCode int, JSONBody []byte) events.APIGatewayProxyResponse {
	var buf bytes.Buffer
	json.HTMLEscape(&buf, JSONBody)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	return resp
}
