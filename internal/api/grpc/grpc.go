package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/api/grpc/handlers"
	"github.com/chains-lab/places-svc/internal/api/grpc/interceptors"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logium.Logger, app *app.App) error {
	requestId := interceptors.RequestID()
	auth := interceptors.Auth(cfg.JWT.User.AccessToken.SecretKey, cfg.JWT.Service.SecretKey)

	grpcUserServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestId,
			auth,
		),
	)

	service := handlers.NewService(app, cfg, log)

	Pproto.RegisterDistributorServiceServer(grpcUserServer, service)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcUserServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gRPC server …")
		grpcUserServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
