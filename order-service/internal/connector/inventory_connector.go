package connector

import (
	"context"

	"github.com/ymanshur/synasishouse/pb"
	"google.golang.org/grpc"
)

type Inventorier interface {
	CheckStock(ctx context.Context, arg CheckStockParams) (bool, error)
}

type inventoryConn struct {
	client pb.InventoryClient
}

func NewInventory(cc *grpc.ClientConn) Inventorier {
	return &inventoryConn{
		client: pb.NewInventoryClient(cc),
	}
}

type CheckStockParams struct {
	Code   string
	Amount int32
}

func (c *inventoryConn) CheckStock(ctx context.Context, arg CheckStockParams) (bool, error) {
	res, err := c.client.CheckStock(ctx, &pb.CreateStockRequest{
		Code:   arg.Code,
		Amount: arg.Amount,
	})
	if err != nil {
		return false, err
	}

	return res.GetIsAvailable(), nil
}
