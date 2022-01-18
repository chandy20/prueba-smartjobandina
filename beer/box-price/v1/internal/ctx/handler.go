package ctx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

//htpClientInterface contract for http client
type htpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type responseLambda struct {
	PriceTotal float64 `json:"price_total"`
}

//Handler main struct for lambda
type Handler struct {
	beersRepository   beerRepositoryInterface
	httpClient        htpClientInterface
	accessKeyCurrency string //5d50da19b1c34e648b2aca13e2ac05fc
	logger            *logrus.Logger
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

	currency := req.QueryStringParameters["currency"]
	if strings.TrimSpace(currency) == "" {
		return lib.ResponseError(http.StatusBadRequest, errors.New("currency_can_not_be_empty")), nil

	}

	quantityString := req.QueryStringParameters["quantity"]
	if strings.TrimSpace(quantityString) == "" {
		return lib.ResponseError(http.StatusBadRequest, errors.New("quantity_can_not_be_empty")), nil

	}

	ID, err := strconv.Atoi(IDString)
	if err != nil {
		return lib.ResponseError(http.StatusBadRequest, errors.New("beerID_is_not_a_number")), nil

	}
	_, err = strconv.Atoi(quantityString)
	if err != nil {
		return lib.ResponseError(http.StatusBadRequest, errors.New("quantity_is_not_a_number")), nil

	}

	beer, err := h.beersRepository.Find(ID)
	if err != nil {
		return lib.ResponseError(http.StatusInternalServerError, err), nil

	}

	if beer.ID == 0 {
		return lib.ResponseError(http.StatusNotFound, errors.New("beerID_does_not_exist")), nil
	}
	uri := fmt.Sprintf(
		"https://api.currencylayer.com/convert?access_key=%v&from=%v&to=%v&amount=%v",
		h.accessKeyCurrency,
		beer.Currency,
		currency,
		beer.Price,
	)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		body, _ := json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		return lib.JSONResponse(http.StatusInternalServerError, body), nil
	}

	response, err := h.httpClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		body, _ := json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		return lib.JSONResponse(http.StatusBadRequest, body), nil
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		body, _ := json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		return lib.JSONResponse(http.StatusBadRequest, body), nil
	}
	h.logger.WithField("response_body", responseBody).Info("response from api")

	//Queda desarrollo inconcludo por que libreria suministrada para la prueba es de pago
	//faltó poner api key en variables de entorno e inyectarla al repositorio
	//Hacer operacion matematica con la respuesta y la cantidad de cervezas por default (6) y retornar respuesta exitosa

	return lib.EmptyResponse(http.StatusOK), nil

}

//NewHandler construct for Handler
func NewHandler(
	beersRepository beerRepositoryInterface,
	httpClient htpClientInterface,
	accessKeyCurrency string,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		beersRepository:   beersRepository,
		httpClient:        httpClient,
		accessKeyCurrency: accessKeyCurrency,
		logger:            logger,
	}
}
