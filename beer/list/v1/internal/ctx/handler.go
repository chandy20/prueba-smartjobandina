package ctx

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
)

//beerRepositoryInterface contract for beer repository
type beerRepositoryInterface interface {
	List() ([]model.Beer, error)
}

//Handler main struct for lambda
type Handler struct {
	beersRepository beerRepositoryInterface
	logger          *logrus.Logger
}

//Handler main function for lambda
func (h *Handler) Handler(
	_ context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	logger := h.logger.WithField("request_body", req.Body)
	logger.Info("Beginning of execution of lambda")

	beers, err := h.beersRepository.List()
	if err != nil {
		logger.WithError(err).Error("error finding beers")
		return responseError(http.StatusInternalServerError, err), nil
	}
	if len(beers) == 0 {
		return emptyResponse(http.StatusAccepted), nil
	}

	response, err := json.Marshal(beers)
	if err != nil {
		logger.WithError(err).Error("error marshaling beers")
		return responseError(http.StatusInternalServerError, err), nil
	}

	return jsonResponse(http.StatusOK, response), nil
}

// emptyResponse func for response
func emptyResponse(statusCode int) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
	return resp
}

//responseError function to build json error
func responseError(statusCode int, err error) events.APIGatewayProxyResponse {

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

// jsonResponse receives a JSON body and a code and returns a Response of type APIGatewayProxyResponse
func jsonResponse(statusCode int, JSONBody []byte) events.APIGatewayProxyResponse {
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

//NewHandler construct for Handler
func NewHandler(
	beersRepository beerRepositoryInterface,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		beersRepository: beersRepository,
		logger:          logger,
	}
}
