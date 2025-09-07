package handlers

import (
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
)

type Adapter struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewAdapter(cfg config.Config, log logium.Logger, a *app.App) Adapter {
	return Adapter{
		app: a,
		log: log,
		cfg: cfg,
	}
}

func (a Adapter) Log(r *http.Request) logium.Logger {
	return a.log
}
