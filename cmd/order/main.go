package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/tgkzz/order/config"
	"github.com/tgkzz/order/internal/app"
	pkgLogger "github.com/tgkzz/order/pkg/logger"
)

const (
	cfgPath = "CONFIG_PATH"
	env     = "ENV"
)

func main() {
	cPath := os.Getenv(cfgPath)
	cfg := config.MustRead(cPath)

	var logger *slog.Logger
	switch cfg.Env != "" {
	case true:
		logger = pkgLogger.SetupLogger(env)
	default:
		logger = pkgLogger.SetupLogger("local")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered Error after panic", slog.Any("panic", r))
		}
	}()

	c := context.Background()

	ctx, stop := signal.NotifyContext(c, os.Interrupt)
	defer stop()

	a, err := app.New(*cfg, logger)
	if err != nil {
		panic(err)
	}

	go func() {
		a.HttpServer.MustRun()
	}()

	go func() {
		a.GrpcServer.MustRun()
	}()

	<-ctx.Done()
	a.GrpcServer.Stop()
	if err = a.HttpServer.Stop(ctx); err != nil {
		panic(err)
	}
}
