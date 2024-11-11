package httpHandlers

import (
	v1 "github.com/Cheasezz/anSpace/backend/internal/transport/http/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Cheasezz/anSpace/backend/docs"
)

type Handlers struct {
	v1 *v1.Handlers
}

func NewHandlers(d v1.Deps) *Handlers {
	return &Handlers{
		v1: v1.NewHandlers(d),
	}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), corsMiddleware(h.v1.Config.FrontendOrigins))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, v1.ErrorResponse{Message: "Page not found"})
	})

	api := router.Group("/api")
	{
		h.v1.InitRoutes(api)
	}
	return router
}
