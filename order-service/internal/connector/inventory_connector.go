package connector

import (
	"context"

	"github.com/ymanshur/synasishouse/pb"
	"google.golang.org/grpc"
)

type Inventorier interface {
	CheckStock(ctx context.Context, arg StockParams) (bool, error)
}

type inventoryConn struct {
	client pb.InventoryClient
}

func NewInventory(cc *grpc.ClientConn) Inventorier {
	return &inventoryConn{
		client: pb.NewInventoryClient(cc),
	}
}

type StockParams struct {
	Code   string
	Amount int32
}

func (c *inventoryConn) CheckStock(ctx context.Context, arg StockParams) (bool, error) {
	res, err := c.client.CheckStock(ctx, &pb.GetStockRequest{
		Code:   arg.Code,
		Amount: arg.Amount,
	})
	if err != nil {
		return false, err
	}

	return res.GetIsAvailable(), nil
}
