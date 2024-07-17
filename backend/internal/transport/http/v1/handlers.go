package v1

import (
	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	*Auth
}

type Deps struct {
	Services     *service.Services
	TokenManager auth.TokenManager
	ConfigHTTP   config.HTTP
	Log          logger.Logger
}

func NewHandlers(d Deps) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(d),
	}
}

func (h *Handlers) InitRoutes(router *gin.RouterGroup) {

	v1 := router.Group("/v1")
	{
		h.initAuthRoutes(v1)
		// wel := v1.Group("/wel", h.userIdentity)
		// {
		// 	wel.GET("/come", h.Welcome)
		// }
	}
}

// func (h *Handlers) Welcome(c *gin.Context) {
// 	id, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
// 		return
// 	}

// 	logrus.Printf("hello Its Welcome Func!!! And userID is: %s", id)
// 	c.JSON(http.StatusOK, gin.H{
// 		"cho": "Norm vse",
// 	})
// }
