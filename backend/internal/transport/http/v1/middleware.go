package v1

import (
	"fmt"
	"net/http"
	"strings"

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

// Middleware for identify user with auth header
func (h *Auth) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, h.log, http.StatusUnauthorized, errEmptyAuthHeader)
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, h.log, http.StatusUnauthorized, errInvalidAuthHeader)
		return
	}

	userId, err := h.TokenManager.Parse(headerParts[1])
	if err != nil {
		if err.Error() == errAccessTokenIsExpired.Error() {
			newErrorResponse(c, h.log, http.StatusUnauthorized, err)
			return
		}
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	c.Set(userCtx, userId)
}

// This function return user id from gin context
func getUserIdFrmCtx(c *gin.Context) (uuid.UUID, error) {
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
