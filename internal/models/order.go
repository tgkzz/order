package models

import (
	order1 "github.com/tgkzz/order/gen/go/order"
)

type Order struct {
	Id         string  `json:"id" bson:"_id,omitempty"`
	Username   string  `json:"username" bson:"username"`
	TotalPrice float64 `json:"total_price" bson:"total_price"`
	Items      []Item  `json:"items" bson:"items"`
}

type Item struct {
	ItemId   string  `json:"item_id" bson:"item_id"`
	Name     string  `json:"name" bson:"name"`
	Price    float64 `json:"price" bson:"price"`
	Currency string  `json:"currency" bson:"currency"`
}

func FromDtoItemToItem(req []*order1.CreateOrderItemRequest) []Item {
	res := make([]Item, len(req))
	for i, item := range req {
		res[i] = Item{
			Name:     item.Name,
			Price:    float64(item.Price),
			Currency: item.Currency,
		}
	}
	return res
}

func FromItemToOrderResponseItem(req []Item) []*order1.GetOrderItemRequest {
	res := make([]*order1.GetOrderItemRequest, len(req))

	for i, item := range req {
		res[i] = &order1.GetOrderItemRequest{
			ItemId:   item.ItemId,
			Name:     item.Name,
			Price:    float32(item.Price),
			Currency: item.Currency,
		}
	}

	return res
}
