package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/cmd/config"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	server *http.Server
	router *chi.Mux

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger) Service {
	logger := log.WithField("module", "rest")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:              cfg.Server.Port,
		ReadTimeout:       cfg.Server.Timeouts.Read * time.Second,
		ReadHeaderTimeout: cfg.Server.Timeouts.ReadHeader * time.Second,
		WriteTimeout:      cfg.Server.Timeouts.Write * time.Second,
		IdleTimeout:       cfg.Server.Timeouts.Idle * time.Second,
		Handler:           router,
	}

	return Service{
		router: router,
		server: server,
		log:    logger,
		cfg:    cfg,
	}
}

func (s *Service) start(ctx context.Context) {
	go func() {
		s.log.Infof("Starting server on port %s", s.cfg.Server.Port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (s *Service) stop(ctx context.Context) {
	s.log.Info("Shutting down server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Errorf("Server shutdown failed: %v", err)
	}
}
