package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/internal/service"
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
	service service.Auth
	config  config.HTTP
	log     logger.Logger
	mdlwrs  *Middlewares
}

func (h *Auth) initAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
		auth.DELETE("/logout", h.logOut)
		auth.POST("/genpasrcode", h.genPassResetCode)
		auth.POST("/refresh", h.refreshAccessToken)
		auth.GET("/me", h.mdlwrs.userIdentity, h.me)
	}
}

func NewAuthHandler(d Deps, m *Middlewares) *Auth {
	return &Auth{
		service: d.Services.Auth,
		config:  d.ConfigHTTP,
		log:     d.Log,
		mdlwrs:  m,
	}
}

// @Tags auth
// @Summary create account
// @Description create account in db and return access token in JSON and refresh token in cookies
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.AuthCredentials true "signUp input"
// @Success 200 {object} auth.ATknInfo
// @Header 200 {string} Set-Cookie "refreshToken. Example: "RefreshToken=9838c59cff93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None" "
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/v1/auth/signup [post]
func (h *Auth) signUp(c *gin.Context) {
	var input core.AuthCredentials

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateEmail(input.Email); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}
	if err := h.validatePass(input.Password); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.service.SignUp(c, input)
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
// @Param input body core.AuthCredentials true "signin input"
// @Success 200 {object} auth.ATknInfo
// @Header 200 {string} Set-Cookie "refreshToken. Example: "RefreshToken=9838c59cff93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None" "
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/v1/auth/signin [post]
func (h *Auth) signIn(c *gin.Context) {
	var input core.AuthCredentials

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateEmail(input.Email); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}
	if err := h.validatePass(input.Password); err != nil {
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
// @Param Cookie header string true "refresh token in cookies"
// @Success 200 {object} auth.ATknInfo "response has emty accessToken"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/v1/auth/logout [delete]
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
// @Param Cookie header string true "refresh token in cookies"
// @Success 200 {object} auth.ATknInfo
// @Header 200 {string} Set-Cookie "refreshToken. Example: "RefreshToken=9838c59cff93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None" "
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
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

	c.JSON(http.StatusOK, tokens.Access)
}

// @Tags auth
// @Summary return curent username
// @Description return curent username
// @ID me
// @Produce  json
// @Success 200 {object} userResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Security		bearerAuth
// @Router /api/v1/auth/me [get]
func (h *Auth) me(c *gin.Context) {
	usrId, err := h.mdlwrs.getUserIdFrmCtx(c)
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
// @Param email body core.Email true "email input"
// @Produce  json
// @Success 200 "password reset code saved in db and sent on specified email"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/v1/genpasrcode [post]
func (h *Auth) genPassResetCode(c *gin.Context) {
	var email core.Email

	if err := c.ShouldBindJSON(&email); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.validateEmail(email.Email); err != nil {
		newErrorResponse(c, h.log, http.StatusBadRequest, err)
		return
	}

	if err := h.service.GenPassResetCode(c, email.Email); err != nil {
		newErrorResponse(c, h.log, http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Auth) validateEmail(email string) error {
	var (
		trimE = strings.TrimSpace(email)
	)

	if len(trimE) == 0 {
		return errEmptyEmailOrPass
	}

	if ok := govalidator.IsExistingEmail(trimE); !ok {
		return errIncorrectEmail
	}
	return nil
}

func (h *Auth) validatePass(pass string) error {
	var (
		trimP = strings.TrimSpace(pass)
	)

	if len(trimP) == 0 {
		return errEmptyEmailOrPass
	}

	if ok := govalidator.MinStringLength(trimP, "12"); !ok {
		return errShortPass
	}

	return nil
}
