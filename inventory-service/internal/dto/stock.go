package dto

import (
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	pb "github.com/ymanshur/synasishouse/proto"
)

func StockPBToRequest(i *pb.StockRequest) presentation.StockRequest {
	return presentation.StockRequest{
		ProductCode: i.GetProductCode(),
		Amount:      i.GetAmount(),
	}
}

func StocksPBToRequest(items []*pb.StockRequest) []presentation.StockRequest {
	p := []presentation.StockRequest{}
	for _, i := range items {
		p = append(p, StockPBToRequest(i))
	}
	return p
}
