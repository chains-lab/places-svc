package middlewares

import (
	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
)

type Service struct {
	log logium.Logger
}

func New(cfg internal.Config, log logium.Logger) Service {
	return Service{
		log: log,
	}
}
