package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/sirupsen/logrus"
)

func Run(args []string) bool {
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	log := logium.NewLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	log.Info("Starting server...")

	var (
		service = kingpin.New("chains-auth", "")
		runCmd  = service.Command("run", "run command")

		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
		return false
	}

	var wg sync.WaitGroup

	cmd, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		StartServices(ctx, cfg, log, &wg, &application)
	case migrateUpCmd.FullCommand():
		err = dbx.MigrateUp(cfg)
	case migrateDownCmd.FullCommand():
		err = dbx.MigrateDown(cfg)
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All services stopped")
	}

	return true
}
