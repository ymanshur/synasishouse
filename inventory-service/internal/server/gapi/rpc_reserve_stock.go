package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/dto"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	pb "github.com/ymanshur/synasishouse/proto"
)

func (r *Server) ReserveStock(ctx context.Context, req *pb.CreateStockRequest) (*pb.StockResponse, error) {
	res, err := r.stockUseCase.Reserve(ctx, presentation.CreateStockRequest{
		Stocks: dto.StocksPBToRequest(req.Stocks),
	})
	if err != nil {
		return nil, translationError(err)
	}

	return &pb.StockResponse{
		IsAvailable: res.IsAvailable,
	}, nil
}
