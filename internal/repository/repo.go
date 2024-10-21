package repository

import (
	"context"
	"github.com/tgkzz/order/internal/models"
	mongoRepo "github.com/tgkzz/order/internal/repository/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IOrderRepository interface {
	CreateNewOrder(ctx context.Context, order models.Order) (string, error)
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*models.Order, error)
}

const OrderCollection = "order"

func NewMongoOrderRepository(ctx context.Context, uri string) (IOrderRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	conn, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = conn.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err = conn.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &mongoRepo.OrderRepository{
		Coll: conn.Database("order").Collection(OrderCollection),
	}, nil
}
