//go:build wireinject
// +build wireinject

package di

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var stdSet = wire.NewSet(
	session.NewSession,
	dynamodb.New,
	logrus.New,
	providerBeerRepository,
	provideNewHandler,
	providerAWSConfig,

	wire.Bind(new(client.ConfigProvider), new(*session.Session)),
)
