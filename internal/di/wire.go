//go:build wireinject
// +build wireinject

package di

import (
	"backendTemp/internal/controller"
	"backendTemp/internal/repository"
	"backendTemp/internal/router"
	"backendTemp/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
)

var ServiceSet = wire.NewSet(
	service.NewAuthService,
)

var ControllerSet = wire.NewSet(
	controller.NewAuthController,
	wire.Struct(new(router.Controllers), "*"),
)

func InitializeApp(db *mongo.Database) (*gin.Engine, error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		ControllerSet,
		router.SetupRouter,
	)
	return &gin.Engine{}, nil
}
