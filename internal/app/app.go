package app

import (
	"log/slog"

	"github.com/tgkzz/order/config"
	"github.com/tgkzz/order/internal/app/grpc"
	appHttp "github.com/tgkzz/order/internal/app/http"
	"github.com/tgkzz/order/internal/service/order"
)

type App struct {
	HttpServer *appHttp.HTTPServer
	GrpcServer *grpc.App
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	orderSrv, err := order.NewOrderService(logger, cfg.Mongo.Uri, cfg.GrpcStorageServer.Host, cfg.GrpcStorageServer.Port)
	if err != nil {
		return nil, err
	}

	httpSrv, err := appHttp.NewHTTPServer("echo", logger, cfg.HttpOrderServer.Port, cfg)
	if err != nil {
		return nil, err
	}

	grpcSrv := grpc.NewApp(logger, orderSrv, cfg.GrpcOrderServer.Port)

	return &App{HttpServer: httpSrv, GrpcServer: grpcSrv}, nil
}
