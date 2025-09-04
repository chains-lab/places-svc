package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
)

func StartServices(ctx context.Context, cfg config.Config, log logium.Logger, wg *sync.WaitGroup, app *app.App) {
	_ = func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}
}
