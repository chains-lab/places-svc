package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/infra/geo"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/chains-lab/places-svc/internal/domain/services/timetable"
	"github.com/chains-lab/places-svc/internal/rest"
	"github.com/chains-lab/places-svc/internal/rest/controller"
	"github.com/chains-lab/places-svc/internal/rest/middlewares"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.New(pg)

	geoGuesser := geo.NewGuesser()

	classSvc := class.NewService(database)
	placeSvc := place.NewService(database, geoGuesser)
	pLocalesSvc := plocale.NewService(database)
	timetableSvc := timetable.NewService(database)

	ctrl := controller.New(cfg, log, classSvc, placeSvc, pLocalesSvc, timetableSvc)

	mdlv := middlewares.New(cfg, log)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })

}
