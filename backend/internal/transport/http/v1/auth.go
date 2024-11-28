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
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

var (
	errEmptyEmailOrPass = fmt.Errorf("all fields must be completed")
	errShortPass        = fmt.Errorf("password must be more than 11 characters")
	errIncorrectEmail   = fmt.Errorf("incorrect email")
)

type Auth struct {
	service      service.Auth
	TokenManager auth.TokenManager
	Config       config.HTTP
	log          logger.Logger
}

func (h *Auth) initAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.GET("/logout", h.logOut)
		auth.POST("/genpassresetcode", h.genPassResetCode)
		auth.POST("/refresh", h.refreshAccessToken)
		auth.GET("/me", h.userIdentity, h.me)
	}
}

func NewAuthHandler(d Deps) *Auth {
	return &Auth{
		service:      d.Services.Auth,
		TokenManager: d.TokenManager,
		Config:       d.ConfigHTTP,
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
	var input core.AuthCredentials

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateEmailAndPass(input.Email, input.Password); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.service.SignUp(c, input)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	newTokenResponse(c, tokens, h.Config)
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
	var input core.AuthCredentials

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateEmailAndPass(input.Email, input.Password); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.service.SignIn(c, input)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	newTokenResponse(c, tokens, h.Config)
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

	newTokenResponse(c, tkns, h.Config)
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
		newTokenResponse(c, tokens, h.Config)
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
		newErrorResponse(c, h.log, http.StatusUnauthorized, err)
		return
	}
	user, err := h.service.GetUser(c, usrId)
	if err != nil {
		newErrorResponse(c, h.log, http.StatusUnauthorized, err)
		return
	}
	c.JSON(http.StatusOK, userResponse{
		User: user,
	})
}

// @Tags auth
// @Summary generate password reset code
// @Description generate and save password reset code into db. Sends code to email
// @ID gen_pass_reset_code
// @Produce  json
// @Success 200 
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/reset [post]
func (h *Auth) genPassResetCode(c *gin.Context) {

	type Email struct {
		Email string `json:"email"`
	}
	var email Email

	if err := c.ShouldBindJSON(&email); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.service.GenPassResetCode(c, email.Email); err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// TODO: Separate validation for email and pass
func (h *Auth) validateEmailAndPass(e string, p string) error {
	var (
		trimE = strings.TrimSpace(e)
		trimP = strings.TrimSpace(p)
	)

	if len(trimE) == 0 || len(trimP) == 0 {
		return errEmptyEmailOrPass
	}

	if ok := govalidator.MinStringLength(trimP, "12"); !ok {
		return errShortPass
	}
	if ok := govalidator.IsExistingEmail(trimE); !ok {
		return errIncorrectEmail
	}
	return nil
}
