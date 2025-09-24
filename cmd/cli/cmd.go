package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/cmd/config"
	"github.com/chains-lab/places-svc/internal/api"
	"github.com/chains-lab/places-svc/internal/api/rest/controller"
	"github.com/chains-lab/places-svc/internal/data/fabric"
	"github.com/chains-lab/places-svc/internal/domain/modules/class"
	"github.com/chains-lab/places-svc/internal/domain/modules/place"
)

func StartServices(ctx context.Context, cfg config.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	database := fabric.NewDatabase(cfg.Database.SQL.URL)
	classMod := class.NewModule(database)
	placeMod := place.NewModule(database)

	Api := api.NewAPI(cfg, log)

	run(func() {
		handl := controller.NewService(cfg, log, classMod, placeMod)

		Api.RunRest(ctx, handl)
	})
}
