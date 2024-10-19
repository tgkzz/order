package mongo

import (
	"context"

	"github.com/tgkzz/order/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	Coll *mongo.Collection
}

func (or *OrderRepository) CreateNewOrder(ctx context.Context, order models.Order) (string, error) {
	return "", nil
}

func (or *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	return nil
}

func (or *OrderRepository) GetOrderById(ctx context.Context, id string) (*models.Order, error) {
	return nil, nil
}
