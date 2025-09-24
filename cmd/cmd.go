package cmd

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
	"github.com/chains-lab/places-svc/internal/api"
	"github.com/chains-lab/places-svc/internal/api/rest/controller"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	database := data.NewDatabase(cfg.Database.SQL.URL)
	classMod := class.NewService(database)
	placeMod := place.NewService(database)

	Api := api.NewAPI(cfg, log)

	run(func() {
		handl := controller.NewService(cfg, log, classMod, placeMod)

		Api.RunRest(ctx, handl)
	})
}
