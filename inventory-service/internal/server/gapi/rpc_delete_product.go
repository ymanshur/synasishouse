package gapi

import (
	"context"

	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Server) DeleteProduct(ctx context.Context, req *pb.GetProductRequest) (*emptypb.Empty, error) {
	err := r.productUseCase.Delete(ctx, presentation.GetProductRequest{
		ID: req.GetId(),
	})
	return nil, translationError(err)
}
