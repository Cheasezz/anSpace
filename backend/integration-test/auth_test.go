package integration_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func (s *APITestSuite) TestSignUp() {
	_, err := s.db.Pool.Exec(context.Background(), "truncate users, users_sessions")
	if err != nil {
		s.logger.Error("db exec error: %s", err.Error())
	}

	var email string
	inputBody := `{"Email": "Cheasezz@gmail.com", "Password": "qwerty123456"}`
	r := s.Require()

	resp, err := http.Post("http://"+s.server.HttpServer.Addr+"/api/v1/auth/signup", "json", bytes.NewBufferString(inputBody))
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	err = s.db.Scany.Get(context.Background(), s.db.Pool, &email, `select email from users where email='Cheasezz@gmail.com'`)
	if err != nil {
		s.logger.Error("FromTestSignUp db scany get error: %s", err.Error())
	}
	var st []byte
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal("Cheasezz@gmail.com", email)
	r.Contains(string(st), `{"accessToken":`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")

}

func (s *APITestSuite) TestSignIn() {
	inputSignIn := `{"Email": "Cheasezz@gmail.com", "Password": "qwerty123456"}`
	r := s.Require()

	resp, err := http.Post("http://"+s.server.HttpServer.Addr+"/api/v1/auth/signin", "json", bytes.NewBufferString(inputSignIn))
	if err != nil {
		s.logger.Error("http post error: %s", err.Error())
	}
	var st []byte
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(string(st), `{"accessToken":`)
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

	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal(string(st), `{"accessToken":""}`)
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

	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(string(st), `{"accessToken":"`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")
	r.NotEmpty(resp.Cookies()[0].Value)
	r.NotEqual(resp.Cookies()[0].Value, s.userCookie)
}

func (s *APITestSuite) TestMe() {
	var st []byte

	r := s.Require()

	req, err := http.NewRequest("GET", "http://"+s.server.HttpServer.Addr+"/api/v1/auth/me", nil)
	if err != nil {
		s.logger.Error("http get error: %s", err.Error())
	}

	req.Header.Add("Authorization", s.accessToken)
	req.AddCookie(&http.Cookie{
		Name:  "RefreshToken",
		Value: s.userCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	st, _ = io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(string(st), `{"user":{"email":"Cheasezz@gmail.com",`)
}
