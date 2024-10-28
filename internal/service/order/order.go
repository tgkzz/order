package order

import (
	"context"
	"errors"
	"log/slog"

	repoErrs "github.com/tgkzz/order/internal/repository/erros"
	"github.com/tgkzz/order/pkg/grpc/storage"

	"github.com/tgkzz/order/internal/models"
	"github.com/tgkzz/order/internal/repository"
	"github.com/tgkzz/order/pkg/logger"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order models.Order) (string, error)
	GetOrder(ctx context.Context, id string) (models.Order, error)
}

type orderService struct {
	logger          *slog.Logger
	orderRepository repository.IOrderRepository
	storageClient   storage.StorageClient
}

var ErrOrderNotFound = errors.New("order not found")

func NewOrderService(logger *slog.Logger, mongoDbUri, storageHost, storagePort string) (OrderService, error) {
	repo, err := repository.NewMongoOrderRepository(context.TODO(), mongoDbUri)
	if err != nil {
		return nil, err
	}

	storageCli, err := storage.NewStorageClient(storageHost, storagePort, logger)
	if err != nil {
		return nil, err
	}

	return &orderService{
		logger:          logger,
		orderRepository: repo,
		storageClient:   storageCli,
	}, nil
}

func (or *orderService) CreateOrder(ctx context.Context, order models.Order) (string, error) {
	const op = "orderService.CreateOrder"

	log := or.logger.With(
		slog.String("op", op),
		slog.Any("order", order),
	)

	// we need to check availability of order
	if err := or.storageClient.CreateOrder(ctx, order.Username, order.Items); err != nil {
		log.Error("error while creating order in storage service", slog.String("err", err.Error()))
		return "", err
	}

	// also may check how money does user have

	// and may even add auth

	id, err := or.orderRepository.CreateNewOrder(ctx, order)
	if err != nil {
		log.Error("failed to create order", logger.Err(err))
		return "", err
	}

	return id, nil
}

func (or *orderService) GetOrder(ctx context.Context, id string) (models.Order, error) {
	const op = "orderService.GetOrder"

	log := or.logger.With(
		slog.String("op", op),
		slog.String("id", id),
	)

	res, err := or.orderRepository.GetOrderById(ctx, id)
	if err != nil {
		log.Error("failed to get order", logger.Err(err))
		if errors.Is(err, repoErrs.ErrNotFound) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, err
	}

	return *res, nil
}
