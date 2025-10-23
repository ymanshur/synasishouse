package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
)

func (r *Server) ReserveStock(ctx context.Context, req *pb.CreateStockRequest) (*pb.StockResponse, error) {
	res, err := r.stockUseCase.Reserve(ctx, presentation.CreateStockRequest{
		Code:   req.GetCode(),
		Amount: req.GetAmount(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	return &pb.StockResponse{
		IsAvailable: res.IsAvailable,
	}, nil
}
