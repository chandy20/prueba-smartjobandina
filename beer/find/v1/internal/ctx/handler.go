package ctx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/chandy20/prueba-smartjobandina/beer/lib"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
)

//beerRepositoryInterface contract for beer repository
type beerRepositoryInterface interface {
	Find(int) (model.Beer, error)
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
	IDString := req.PathParameters["beerID"]
	if strings.TrimSpace(IDString) == "" {
		return lib.ResponseError(http.StatusBadRequest, errors.New("beerID_can_not_be_empty")), nil
	}

	ID, err := strconv.Atoi(IDString)
	if err != nil {
		return lib.ResponseError(http.StatusBadRequest, errors.New("beerID_is_not_a_number")), nil

	}

	beer, err := h.beersRepository.Find(ID)
	if err != nil {
		return lib.ResponseError(http.StatusInternalServerError, err), nil

	}

	if beer.ID == 0 {
		return lib.ResponseError(http.StatusNotFound, errors.New("beerID_does_not_exist")), nil
	}

	response, err := json.Marshal(beer)
	if err != nil {
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
