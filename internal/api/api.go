package api

import (
	"context"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/cmd/config"
	"github.com/chains-lab/places-svc/internal/api/rest"
	"github.com/chains-lab/places-svc/internal/api/rest/controller"
)

type API struct {
	Rest rest.Service
}

func NewAPI(cfg config.Config, log logium.Logger) API {
	return API{
		Rest: rest.NewRest(cfg, log),
	}
}

func (a *API) RunRest(ctx context.Context, con controller.Service) {
	a.Rest.Run(ctx, con)
}
