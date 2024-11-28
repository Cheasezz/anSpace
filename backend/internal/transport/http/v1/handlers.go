package v1

import (
	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Config config.HTTP
	*Middlewares
	*Auth
}

type Deps struct {
	Services     *service.Services
	TokenManager auth.TokenManager
	ConfigHTTP   config.HTTP
	Log          logger.Logger
}

func NewHandlers(d Deps) *Handlers {
	mdlwrs := NewMiddlewares(d)
	return &Handlers{
		Config:      d.ConfigHTTP,
		Middlewares: mdlwrs,
		Auth:        NewAuthHandler(d, mdlwrs),
	}
}

func (h *Handlers) InitRoutes(router *gin.RouterGroup) {

	v1 := router.Group("/v1")
	{
		h.initAuthRoutes(v1)
	}
}
