package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
)

func (r *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := r.productUseCase.Delete(ctx, presentation.GetProductRequest{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	res := pb.DeleteProductResponse{
		IsSuccess: true,
	}
	return &res, nil
}
