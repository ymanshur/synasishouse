package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
)

func (r *Server) ReserveStock(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	err := r.stockUseCase.Check(ctx, presentation.StockRequest{
		Code:   req.GetCode(),
		Amount: req.GetAmount(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	res := pb.StockResponse{
		IsSuccess: true,
	}
	return &res, nil
}
