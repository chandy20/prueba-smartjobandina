package ctx

import (
	"context"
	"encoding/json"
	"github.com/chandy20/prueba-smartjobandina/beer/lib"
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
		return lib.ResponseError(http.StatusInternalServerError, err), nil
	}
	if len(beers) == 0 {
		return lib.EmptyResponse(http.StatusAccepted), nil
	}

	response, err := json.Marshal(beers)
	if err != nil {
		logger.WithError(err).Error("error marshaling beers")
		return lib.ResponseError(http.StatusInternalServerError, err), nil
	}

	return lib.JSONResponse(http.StatusOK, response), nil
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
