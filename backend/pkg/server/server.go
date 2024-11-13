package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/Cheasezz/anSpace/backend/config"
)

const shutdownTimeout = time.Second * 3

type Server struct {
	HttpServer      *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(cfg config.HTTP, handler http.Handler) *Server {

	s := &Server{
		HttpServer: &http.Server{
			Addr:           cfg.Host + ":" + cfg.Port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
		notify:          make(chan error),
		shutdownTimeout: shutdownTimeout,
	}

	s.Run()

	return s
}

func (s *Server) Run() {
	go func() {
		s.notify <- s.HttpServer.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.HttpServer.Shutdown(ctx)
}
