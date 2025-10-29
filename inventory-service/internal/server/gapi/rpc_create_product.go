package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	pb "github.com/ymanshur/synasishouse/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product, err := r.productUseCase.Create(ctx, presentation.CreateProductRequest{
		Code:  req.GetCode(),
		Total: req.GetTotal(),
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
