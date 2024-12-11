//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/coderlewin/kratosinit/internal/biz"
	"github.com/coderlewin/kratosinit/internal/conf"
	"github.com/coderlewin/kratosinit/internal/data"
	"github.com/coderlewin/kratosinit/internal/pkg/auth"
	"github.com/coderlewin/kratosinit/internal/pkg/middleware"
	"github.com/coderlewin/kratosinit/internal/server"
	"github.com/coderlewin/kratosinit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Jwt, log.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			data.ProviderSet,
			biz.ProviderSet,
			service.ProviderSet,
			middleware.ProviderSet,
			auth.ProviderSet,
			newApp,
		),
	)
}
