package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

var (
	errEmptyLoginOrPass = fmt.Errorf("all fields must be completed")
	errShortPass        = fmt.Errorf("password must be more than 11 characters")
)

type Auth struct {
	service      service.Auth
	TokenManager auth.TokenManager
	config       config.HTTP
	log          logger.Logger
}

func (h *Auth) initAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.GET("/logout", h.logOut)
		auth.POST("/refresh", h.refreshAccessToken)
		auth.GET("/me", h.userIdentity, h.me)
	}
}

func NewAuthHandler(d Deps) *Auth {
	return &Auth{
		service:      d.Services.Auth,
		TokenManager: d.TokenManager,
		config:       d.ConfigHTTP,
		log:          d.Log,
	}
}

// @Tags auth
// @Summary create account
// @Description create account in data base and return access token in JSON and refresh token in cookies
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.SignUp true "signUp input"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/auth/sign-up [post]
func (h *Auth) signUp(c *gin.Context) {
	var input core.SignUp

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateLoginAndPass(input.Username, input.Password); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	_, err := h.service.SignUp(c, input)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	inputSignIn := core.SignIn{Username: input.Username, Password: input.Password}
	tokens, err := h.service.SignIn(c, inputSignIn)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	newTokenResponse(c, tokens, h.config)
}

// @Tags auth
// @Summary login to account
// @Description login to accont with username and password and return access token in JSON and refresh token in cookies
// @ID login-to-account
// @Accept  json
// @Produce  json
// @Param input body core.SignIn true "signin input"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/auth/sign-in [post]
func (h *Auth) signIn(c *gin.Context) {
	var input core.SignIn

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.service.SignIn(c, input)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	newTokenResponse(c, tokens, h.config)
}

// @Tags auth
// @Summary logout
// @Description accept refresh token from cookie, and return empty tokens
// @ID logout
// @Produce  json
// @Success 200 {object} tokenResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/auth/logout [get]
func (h *Auth) logOut(c *gin.Context) {
	rt, err := c.Cookie(rtCookieName)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusUnauthorized, err)
		return
	}

	tkns, err := h.service.LogOut(c, rt)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	newTokenResponse(c, tkns, h.config)
}

// @Tags auth
// @Summary refresh access token
// @Description accept refresh token from cookie, and return new access token
// @ID refresh-access-token
// @Produce  json
// @Success 200 {object} tokenResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Auth) refreshAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie(rtCookieName)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusUnauthorized, err)
		return
	}

	tokens, err := h.service.RefreshAccessToken(c, refreshToken)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	if tokens.Refresh.Token != "" {
		newTokenResponse(c, tokens, h.config)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Access: tokens.Access,
	})
}

// @Tags auth
// @Summary return curent username
// @Description return curent username
// @ID me
// @Produce  json
// @Success 200 {object} userResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/me [get]
func (h *Auth) me(c *gin.Context) {
	usrId, err := getUserIdFrmCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{
			Message: err.Error(),
		})
		return
	}
	usrName, err := h.service.GetUser(c, usrId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userResponse{
		User: usrName,
	})
}

func (h *Auth) validateLoginAndPass(l string, p string) error {
	var (
		trimL = strings.TrimSpace(l)
		trimP = strings.TrimSpace(p)
	)

	if len(trimL) == 0 || len(trimP) == 0 {
		return errEmptyLoginOrPass
	}

	if len([]rune(trimP)) < 12 {
		return errShortPass
	}
	return nil
}
