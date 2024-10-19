package mongo

import (
	"context"
	"errors"

	"github.com/tgkzz/order/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/tgkzz/order/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	Coll *mongo.Collection
}

func (or *OrderRepository) CreateNewOrder(ctx context.Context, order models.Order) (string, error) {
	resp, err := or.Coll.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return resp.InsertedID.(string), nil
}

func (or *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	resp, err := or.Coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if resp.DeletedCount == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (or *OrderRepository) GetOrderById(ctx context.Context, id string) (*models.Order, error) {
	filter := bson.M{"_id": id}

	var order models.Order
	if err := or.Coll.FindOne(ctx, filter).Decode(&order); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}
