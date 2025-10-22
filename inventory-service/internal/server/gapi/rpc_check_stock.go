package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
)

func (r *Server) CheckStock(ctx context.Context, req *pb.GetStockRequest) (*pb.StockResponse, error) {
	res, err := r.stockUseCase.Check(ctx, presentation.GetStockRequest{
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
