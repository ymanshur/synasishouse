package connector

import (
	"context"

	pb "github.com/ymanshur/synasishouse/proto"
	"google.golang.org/grpc"
)

type Inventorier interface {
	CheckStock(ctx context.Context, arg CheckStockParams) (bool, error)
	ReserveStock(ctx context.Context, arg ReserveStockParams) (bool, error)
	ReleaseStock(ctx context.Context, arg ReleaseStockParams) (bool, error)
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
	ProductCode string
	Amount      int32
}

type CheckStockParams struct {
	Stocks []StockParams
}

func (c *inventoryConn) CheckStock(ctx context.Context, arg CheckStockParams) (bool, error) {
	in := &pb.CreateStockRequest{
		Stocks: []*pb.StockRequest{},
	}
	for _, stock := range arg.Stocks {
		in.Stocks = append(in.Stocks, &pb.StockRequest{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}
	res, err := c.client.CheckStock(ctx, in)
	if err != nil {
		return false, err
	}

	return res.GetIsAvailable(), nil
}

type ReserveStockParams struct {
	Stocks []StockParams
}

func (c *inventoryConn) ReserveStock(ctx context.Context, arg ReserveStockParams) (bool, error) {
	in := &pb.CreateStockRequest{
		Stocks: []*pb.StockRequest{},
	}
	for _, stock := range arg.Stocks {
		in.Stocks = append(in.Stocks, &pb.StockRequest{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}
	res, err := c.client.ReserveStock(ctx, in)
	if err != nil {
		return false, err
	}

	return res.GetIsAvailable(), nil
}

type ReleaseStockParams struct {
	Stocks []StockParams
}

func (c *inventoryConn) ReleaseStock(ctx context.Context, arg ReleaseStockParams) (bool, error) {
	in := &pb.CreateStockRequest{
		Stocks: []*pb.StockRequest{},
	}
	for _, stock := range arg.Stocks {
		in.Stocks = append(in.Stocks, &pb.StockRequest{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}
	res, err := c.client.ReleaseStock(ctx, in)
	if err != nil {
		return false, err
	}

	return res.GetIsAvailable(), nil
}
