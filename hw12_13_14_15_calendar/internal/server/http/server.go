package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	logg   app.Logger
	// app *app.App
	address string
}

func NewServer(host string, port string, app *app.App, logg app.Logger) *Server {
	api := API{
		App:  app,
		Logg: logg,
	}

	router := mux.NewRouter()
	router.HandleFunc("/hello", api.Hello).Methods("GET")
	router.HandleFunc("/events", api.AddEvent).Methods("POST")
	router.HandleFunc("/events", api.UpdateEvent).Methods("PUT")
	router.HandleFunc("/events", api.DeleteEvent).Methods("DELETE")
	router.HandleFunc("/events", api.ListPerPeriod).Methods("GET")

	server := http.Server{ // nolint: exhaustivestruct
		Addr:         net.JoinHostPort(host, port),
		Handler:      loggingMiddleware(router, logg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: &server,
		logg:   logg,
		// app: app,
		address: net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info(fmt.Sprintf("http server is running on %s", s.address))

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error occurred on attempt to run server: %w", err)
	}

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
