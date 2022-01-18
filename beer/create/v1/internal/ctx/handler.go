package ctx

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/chandy20/prueba-smartjobandina/beer/lib"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

//errBeerAlreadyCreated custom error to represent that a beer is already created
var errBeerAlreadyCreated = errors.New("error_beer_already_created")

//beerRepositoryInterface contract for beer repository
type beerRepositoryInterface interface {
	Find(int) (model.Beer, error)
	Save(model.Beer) error
}

//Handler main struct for lambda
type Handler struct {
	beersRepository beerRepositoryInterface
	logger          *logrus.Logger
}

//go:embed json_files/validation.json
var validationSchema []byte

//Handler main function for lambda
func (h *Handler) Handler(
	_ context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	logger := h.logger.WithField("request_body", req.Body)
	logger.Info("Beginning of execution of lambda")

	var beer model.Beer
	err := json.Unmarshal([]byte(req.Body), &beer)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"message": req.Body,
		}).WithError(err).Error("error while unmarshalling")
		return lib.ResponseError(http.StatusBadRequest, err), nil
	}

	dataToValidate := gojsonschema.NewGoLoader(beer)
	result, err := gojsonschema.Validate(gojsonschema.NewStringLoader(string(validationSchema)), dataToValidate)
	if err != nil {
		logger.WithError(err).Error("error while validating json schema")
		return lib.ResponseError(http.StatusBadRequest, err), nil
	}

	if !result.Valid() {
		logger.WithField("errors", result.Errors()).Error("validation errors found")
		return lib.ResponseError(http.StatusBadRequest, errors.New(result.Errors()[0].String())), nil
	}

	beerGot, err := h.beersRepository.Find(beer.ID)
	if err != nil {
		logger.WithError(err).Error("error finding beer")
		return lib.ResponseError(http.StatusInternalServerError, err), nil
	}

	if beerGot.ID > 0 {
		logger.WithField("beer", beerGot).Error("beer already created")
		return lib.ResponseError(http.StatusConflict, errBeerAlreadyCreated), nil
	}

	err = h.beersRepository.Save(beer)
	if err != nil {
		logger.WithField("beer", beer).Error("error saving  beer")
		return lib.ResponseError(http.StatusInternalServerError, err), nil
	}

	return lib.EmptyResponse(http.StatusCreated), nil
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
