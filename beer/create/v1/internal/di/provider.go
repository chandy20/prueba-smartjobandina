package di

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chandy20/prueba-smartjobandina/beer/create/v1/internal/ctx"
	"github.com/chandy20/prueba-smartjobandina/beer/repository"
	"github.com/sirupsen/logrus"
	"os"
)

func providerAWSConfig() []*aws.Config {
	return nil
}

func providerBeerRepository(
	client *dynamodb.DynamoDB,
	logger *logrus.Logger,
) (*repository.BeerRepository, error) {
	tableBeers := os.Getenv("DYNAMODB_BEERS")
	if tableBeers == "" {
		return nil, errors.New("variable DYNAMODB_BEERS is not defined")
	}
	return repository.NewBeerRepository(client, tableBeers, logger), nil
}

func provideNewHandler(
	beerRepository *repository.BeerRepository,
	logger *logrus.Logger,
) *ctx.Handler {
	return ctx.NewHandler(beerRepository, logger)
}
