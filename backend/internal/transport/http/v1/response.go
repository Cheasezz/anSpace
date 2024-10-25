package v1

import (
	"net/http"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type tokenResponse struct {
	Access string `json:"accessToken"`
}

type userResponse struct {
	User core.User `json:"user"`
}

const rtCookieName = "RefreshToken"

func newErrorResponse(c *gin.Context, l logger.Logger, statusCode int, err error) {
	l.Error(err)
	c.AbortWithStatusJSON(statusCode, errorResponse{err.Error()})
}

func newTokenResponse(c *gin.Context, t auth.Tokens, cfg config.HTTP) {
	c.SetCookie(rtCookieName, t.Refresh.Token, t.Refresh.TTLInSec, "/", cfg.Host, false, true)
	c.JSON(http.StatusOK, tokenResponse{
		Access: t.Access,
	})
}
