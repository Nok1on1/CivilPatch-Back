package router

import (
	_ "backendTemp/docs"
	"backendTemp/internal/controller"
	"backendTemp/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controllers struct {
	AuthController *controller.AuthController
}

// @title           CivilPatch-API
// @version         1.0

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func SetupRouter(c *Controllers) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		public := api.Group("/auth")
		{
			public.POST("/register", c.AuthController.Register)
			public.POST("/login", c.AuthController.Login)
		}
	}

	return r
}
