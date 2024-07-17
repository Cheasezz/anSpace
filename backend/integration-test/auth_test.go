package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (s *APITestSuite) TestSignUp() {
	s.db.Pool.Exec(context.Background(), "truncate users, users_sessions")

	var username string
	inputBody := `{"Name": "Iurii", "Username": "Cheasezz", "Password": "qwerty123456"}`
	r := s.Require()

	resp, err := http.Post("http://"+s.server.HttpServer.Addr+"/api/v1/auth/signup", "json", bytes.NewBufferString(inputBody))
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	err = s.db.Scany.Get(context.Background(), s.db.Pool, &username, `select username from users where username='Cheasezz' and name='Iurii'`)
	if err != nil {
		s.logger.Error("FromTestSignUp db scany get error: %s", err.Error())
	}
	var st []byte
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal("Cheasezz", username)
	r.Contains(fmt.Sprintf("%s", st), `{"accessToken":`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")

}

func (s *APITestSuite) TestSignIn() {
	inputSignIn := `{"Username": "Cheasezz", "Password": "qwerty123456"}`
	r := s.Require()

	resp, err := http.Post("http://"+s.server.HttpServer.Addr+"/api/v1/auth/signin", "json", bytes.NewBufferString(inputSignIn))
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	var st []byte
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(fmt.Sprintf("%s", st), `{"accessToken":`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")
}

func (s *APITestSuite) TestLogOut() {

	var st []byte

	r := s.Require()

	req, err := http.NewRequest("GET", "http://"+s.server.HttpServer.Addr+"/api/v1/auth/logout", nil)
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	req.AddCookie(&http.Cookie{
		Name:  "RefreshToken",
		Value: s.userCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal(fmt.Sprintf("%s", st), `{"accessToken":""}`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")
	r.Empty(resp.Cookies()[0].Value)
}

// This case refresh both tokens (access and refresh tokens)
func (s *APITestSuite) TestRefreshAccessToken() {
	var st []byte

	r := s.Require()

	req, err := http.NewRequest("POST", "http://"+s.server.HttpServer.Addr+"/api/v1/auth/refresh", nil)
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	req.AddCookie(&http.Cookie{
		Name:  "RefreshToken",
		Value: s.userCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(fmt.Sprintf("%s", st), `{"accessToken":"`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")
	r.NotEmpty(resp.Cookies()[0].Value)
	r.NotEqual(resp.Cookies()[0].Value, s.userCookie)
}
