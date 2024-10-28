package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tgkzz/order/internal/models"
	storage1 "github.com/tgkzz/storage/gen/go/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient interface {
	CreateOrder(ctx context.Context, username string, items []models.Item) error
	CancelOrder(ctx context.Context, username string) error
}

type Storage struct {
	client storage1.StorageClient
	logger *slog.Logger
}

func NewStorageClient(host, port string, logger *slog.Logger) (StorageClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	c := storage1.NewStorageClient(conn)

	return &Storage{client: c, logger: logger}, nil
}

func (s *Storage) CreateOrder(ctx context.Context, username string, items []models.Item) error {
	const op = "grpcStorageService.CreateOrder"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	reqItem := make([]*storage1.Item, len(items))
	for _, item := range items {
		itemId, err := strconv.Atoi(item.ItemId)
		if err != nil {
			return err
		}
		curr, err := strconv.Atoi(item.Currency)
		if err != nil {
			return err
		}

		reqItem = append(reqItem, &storage1.Item{
			Id:       int32(itemId),
			Name:     item.Name,
			Quantity: 0,
			Price: &storage1.Price{
				Currency: int32(curr),
				Price:    float32(item.Price),
			},
		})
	}

	if resp, err := s.client.CreateOrder(ctx, &storage1.CreateOrderRequest{
		Username: username,
		Items:    reqItem,
	}); err != nil {
		log.Error("error while creating order in storage service",
			slog.String("status", resp.GetResponse().GetStatus()),
			slog.String("errorMsg", resp.GetResponse().GetErr().GetMessage()),
		)
		return err
	}

	return nil
}

func (s *Storage) CancelOrder(ctx context.Context, username string) error {
	return nil
}
