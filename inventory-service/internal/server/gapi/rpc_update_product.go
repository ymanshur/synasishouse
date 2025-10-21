package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Server) UpdateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product, err := r.productUseCase.Create(ctx, presentation.CreateProductRequest{
		Code:  req.Code,
		Total: int32(req.Total),
	})
	if err != nil {
		return nil, translationError(err)
	}

	res := pb.ProductResponse{
		Product: &pb.Product{
			Code:      product.Code,
			Total:     product.Total,
			Reserved:  product.Reserved,
			Available: product.Available,
			UpdatedAt: timestamppb.New(product.UpdatedAt),
			CreatedAt: timestamppb.New(product.CreatedAt),
		},
	}
	return &res, nil
}
