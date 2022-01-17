//go:build wireinject
// +build wireinject

package di

import (
	"github.com/chandy20/prueba-smartjobandina/beer/create/v1/internal/ctx"
	"github.com/google/wire"
)

func Initialize() (*ctx.Handler, error) {
	wire.Build(stdSet)

	return &ctx.Handler{}, nil
}
