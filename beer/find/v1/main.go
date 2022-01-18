package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chandy20/prueba-smartjobandina/beer/find/v1/internal/di"
)

func main() {
	handler, err := di.Initialize()
	if err != nil {
		panic("fatal err: " + err.Error())
	}
	lambda.Start(handler.Handler)
}
