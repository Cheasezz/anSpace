package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewJWT(userId string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (RTInfo, error)
	ValidateRefreshToken(expiresAt time.Time) (int, error)
}

type Manager struct {
	signingKey string
	atTTL      time.Duration
	rtTTL      time.Duration
}

type Tokens struct {
	Access  string
	Refresh RTInfo
}

func NewManager(cfg config.TokenManager) (*Manager, error) {
	if cfg.SigningKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{
		signingKey: cfg.SigningKey,
		atTTL:      cfg.AccessTokenTTL,
		rtTTL:      cfg.RefreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.atTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

type RTInfo struct {
	Token     string
	ExpiresAt time.Time
	TTLInSec  int
}

func (m *Manager) NewRefreshToken() (RTInfo, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().UTC().Unix() + rand.Int63())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return RTInfo{}, err
	}

	refreshToken := RTInfo{
		Token:     fmt.Sprintf("%x", b),
		ExpiresAt: time.Now().Add(m.rtTTL).UTC(),
		TTLInSec:  int(m.rtTTL.Seconds()),
	}
	return refreshToken, nil
}

func (m *Manager) ValidateRefreshToken(expiresAt time.Time) (int, error) {
	var rtDayUntilExpire = int(time.Until(expiresAt).Hours()) / 24
	if rtDayUntilExpire <= 0 {
		return 0, fmt.Errorf("refresh token is expired")
	}
	return rtDayUntilExpire, nil
}
