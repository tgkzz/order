package http

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tgkzz/order/config"
	echoHandler "github.com/tgkzz/order/internal/handler/http/echo"
)

type Handler interface {
	Start(port int) error
	Stop(ctx context.Context) error
}

const (
	echo = "echo"
)

func NewHandler(name string, logger *slog.Logger, cfg config.Config) (Handler, error) {
	switch name {
	case echo:
		return echoHandler.NewHttpServer(cfg, logger)
	default:
		return nil, errors.New("unknown handler")
	}
}
