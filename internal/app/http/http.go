package http

import (
	"context"
	"log/slog"

	"github.com/tgkzz/order/config"
	handlerHttp "github.com/tgkzz/order/internal/handler/http"
)

type HTTPServer struct {
	handler handlerHttp.Handler
	port    int
}

func NewHTTPServer(handlerName string, logger *slog.Logger, port int, cfg config.Config) (*HTTPServer, error) {
	handler, err := handlerHttp.NewHandler(handlerName, logger, cfg)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		handler: handler,
		port:    port,
	}, nil
}

func (s *HTTPServer) MustRun() {
	if err := s.run(); err != nil {
		panic(err)
	}
}

func (s *HTTPServer) run() error {
	return s.handler.Start(s.port)
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.handler.Stop(ctx)
}
