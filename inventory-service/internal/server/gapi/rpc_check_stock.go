package gapi

import (
	"context"
	"fmt"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
)

func (r *Server) CheckStock(ctx context.Context, req *pb.CreateStockRequest) (*pb.StockResponse, error) {
	res, err := r.stockUseCase.Check(ctx, presentation.CreateStockRequest{
		Code:   req.GetCode(),
		Amount: req.GetAmount(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	fmt.Println(res)

	return &pb.StockResponse{
		IsAvailable: res.IsAvailable,
	}, nil
}
