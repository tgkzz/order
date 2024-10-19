package order

import (
	"context"
	"time"

	order1 "github.com/tgkzz/order/gen/go/order"
	"github.com/tgkzz/order/internal/models"
	"github.com/tgkzz/order/internal/service/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TODO: add service
type serverApi struct {
	order1.UnimplementedOrderServiceServer
	orderService order.OrderService
}

func Register(gRPCServer *grpc.Server, orderService order.OrderService) {
	order1.RegisterOrderServiceServer(gRPCServer, &serverApi{orderService: orderService})
}

func (s *serverApi) CreateOrder(ctx context.Context, req *order1.CreateOrderRequest) (*order1.CreateOrderResponse, error) {
	if req.Items == nil {
		return nil, status.Error(codes.InvalidArgument, "items is required")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	orderCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	id, err := s.orderService.CreateOrder(orderCtx,
		models.Order{
			Username:   req.GetUsername(),
			TotalPrice: float64(req.GetTotalPrice()),
			Items:      models.FromDtoItemToItem(req.GetItems()),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &order1.CreateOrderResponse{OrderId: id}, nil
}

func (s *serverApi) GetOrderById(ctx context.Context, req *order1.GetOrderRequest) (*order1.GetOrderResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order id is required")
	}

	orderCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := s.orderService.GetOrder(orderCtx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &order1.GetOrderResponse{OrderId: res.Id, Username: res.Username, Items: models.FromItemToOrderResponseItem(res.Items)}, nil
}

func (s *serverApi) DeleteOrderById(context.Context, *order1.DeleteOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}
