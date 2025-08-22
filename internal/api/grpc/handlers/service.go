package handlers

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/api/grpc/meta"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
)

type Service struct {
	app *app.App
	cfg config.Config
	log logium.Logger
}

func NewService(a *app.App, cfg config.Config, log logium.Logger) *Service {
	return &Service{
		app: a,
		cfg: cfg,
		log: log,
	}
}

func (s Service) Log(ctx context.Context) logium.Logger {
	return s.log.WithField("request_id", meta.RequestID(ctx))
}
