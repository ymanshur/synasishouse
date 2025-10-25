package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product, err := r.productUseCase.Update(ctx, presentation.UpdateProductRequest{
		ID:   req.GetId(),
		Code: req.GetCode(),
	})
	if err != nil {
		return nil, translationError(err)
	}

	res := &pb.ProductResponse{
		Product: &pb.Product{
			Id:        product.ID,
			Code:      product.Code,
			Total:     product.Total,
			Hold:      product.Hold,
			UpdatedAt: timestamppb.New(product.UpdatedAt),
			CreatedAt: timestamppb.New(product.CreatedAt),
		},
	}
	return res, nil
}
