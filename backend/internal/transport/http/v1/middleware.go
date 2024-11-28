package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

var (
	errEmptyAuthHeader      = fmt.Errorf("empty auth header")
	errInvalidAuthHeader    = fmt.Errorf("invalid auth header")
	errUserIdNotFound       = fmt.Errorf("user id not found")
	errAccessTokenIsExpired = fmt.Errorf("Token is expired")
)

type Middlewares struct {
	TokenManager auth.TokenManager
	log          logger.Logger
}

func NewMiddlewares(d Deps) *Middlewares {
	return &Middlewares{
		TokenManager: d.TokenManager,
		log:          d.Log,
	}
}

// Middleware for identify user with auth header
func (m *Middlewares) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, m.log, http.StatusUnauthorized, errEmptyAuthHeader)
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, m.log, http.StatusUnauthorized, errInvalidAuthHeader)
		return
	}

	userId, err := m.TokenManager.Parse(headerParts[1])
	if err != nil {
		if err.Error() == errAccessTokenIsExpired.Error() {
			newErrorResponse(c, m.log, http.StatusUnauthorized, err)
			return
		}
		newErrorResponse(c, m.log, http.StatusInternalServerError, err)
		return
	}

	c.Set(userCtx, userId)
}

// This function return user id from gin context
func (m *Middlewares) getUserIdFrmCtx(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return uuid.UUID{}, errUserIdNotFound
	}
	parsedId, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedId, nil
}
