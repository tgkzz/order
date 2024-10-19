package echo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tgkzz/order/config"
	"github.com/tgkzz/order/internal/service/order"
	"golang.org/x/time/rate"
)

type HttpServer struct {
	orderService order.OrderService
	logger       *slog.Logger
	echoInstance *echo.Echo
}

func NewHttpServer(cfg config.Config, logger *slog.Logger) (*HttpServer, error) {
	orderService, err := order.NewOrderService(logger, cfg.Mongo.Uri)
	if err != nil {
		return nil, err
	}
	return &HttpServer{
		orderService: orderService,
		logger:       logger,
	}, nil
}

func (s *HttpServer) Start(port int) error {
	e := s.routes()

	s.echoInstance = e

	return e.Start(fmt.Sprintf(":%d", port))
}

func (s *HttpServer) Stop(ctx context.Context) error {
	if s.echoInstance == nil {
		return errors.New("echo instance not initialized")
	}

	if err := s.echoInstance.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *HttpServer) routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:   middleware.DefaultSkipper,
		StackSize: 8 << 10,
		LogLevel:  log.ERROR,
	}))

	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 50, ExpiresIn: 5 * time.Minute},
		),
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		Timeout:      time.Second * 30,
		ErrorMessage: RequestTimeout,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			op := fmt.Sprintf("%s", "handler.timeout "+c.Path())
			deadline, _ := c.Request().Context().Deadline()

			l := s.logger.With(
				slog.String("op", op),
				slog.Time("now", time.Now().UTC()),
				slog.Time("deadline", deadline),
			)

			l.Error("request timed out", slog.String("path", c.Path()))
		},
	}))

	e.Use(middleware.Secure())

	v1 := e.Group("/v1")
	{
		order := v1.Group("/order")
		{
			order.POST("/create", s.createOrder)
			order.GET("/:id", s.getOrder)
			order.DELETE("/:id", s.deleteOrder)
		}
	}

	return e
}
