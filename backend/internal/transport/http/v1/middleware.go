package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

var (
	errEmptyAuthHeader   = fmt.Errorf("empty auth header")
	errInvalidAuthHeader = fmt.Errorf("invalid auth header")
	errUserIdNotFound = fmt.Errorf("user id not found")
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
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	c.Set(userCtx, userId)
}

// This function return user id from gin context
func getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errUserIdNotFound
	}

	return id.(string), nil
}
