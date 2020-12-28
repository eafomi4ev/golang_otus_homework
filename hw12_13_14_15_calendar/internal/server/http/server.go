package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	server http.Server
	logg   app.Logger
}

type Handler struct {
}

func (h *Handler) Hello(resp http.ResponseWriter, req *http.Request) {

}

// type Application interface {
// 	// TODO
// }

func New(host string, port string, logg app.Logger) *Server {
	handler := Handler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.Hello)

	server := http.Server{ // nolint: exhaustivestruct
		Addr:         net.JoinHostPort(host, port),
		Handler:      loggingMiddleware(mux, logg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: server, // nolint: govet
		logg:   logg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error occurred on attempt to run server: %w", err)
	}

	<-ctx.Done()
	err = s.Stop(ctx)
	if err != nil {
		return fmt.Errorf("error occurred on attempt to stop server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error on attempt to stop the server: %w", err)
	}
	return nil
}

// TODO
