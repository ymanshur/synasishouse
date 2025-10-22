package usecase

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
	"github.com/ymanshur/synasishouse/inventory/internal/common"
	"github.com/ymanshur/synasishouse/inventory/internal/consts"
	"github.com/ymanshur/synasishouse/inventory/internal/dto"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/inventory/internal/repo"
	"github.com/ymanshur/synasishouse/inventory/internal/typex"
)

type Producter interface {
	Create(ctx context.Context, req presentation.CreateProductRequest) (*presentation.ProductResponse, error)
	Get(ctx context.Context, req presentation.GetProductRequest) (*presentation.ProductResponse, error)
	Update(ctx context.Context, req presentation.UpdateProductRequest) (*presentation.ProductResponse, error)
	Delete(ctx context.Context, req presentation.GetProductRequest) error
}

type productUseCase struct {
	repo repo.Repo
}

func NewProduct(repo repo.Repo) Producter {
	return &productUseCase{repo: repo}
}

func (u *productUseCase) Create(ctx context.Context, req presentation.CreateProductRequest) (*presentation.ProductResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	product, err := u.repo.CreateProduct(ctx, db.CreateProductParams{
		Code:  req.Code,
		Total: req.Total,
	})
	if err != nil {
		if common.PGErrorCode(err) == consts.UniqueViolation {
			return nil, typex.NewUnprocessableEntityError("product unique constraint violated")
		}
		return nil, fmt.Errorf("create product: %w", err)
	}

	res := dto.ProductToResponse(product)
	return &res, nil
}

func (u *productUseCase) Get(ctx context.Context, req presentation.GetProductRequest) (*presentation.ProductResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}

	product, err := u.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, typex.NewNotFoundError("product")
		}
		return nil, fmt.Errorf("get product: %w", err)
	}

	res := dto.ProductToResponse(product)
	return &res, nil
}

func (u *productUseCase) Update(ctx context.Context, req presentation.UpdateProductRequest) (*presentation.ProductResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}

	product, err := u.repo.UpdateProduct(ctx, db.UpdateProductParams{
		Code: req.Code,
		ID:   id,
	})
	if err != nil {
		if common.PGErrorCode(err) == consts.UniqueViolation {
			return nil, typex.NewUnprocessableEntityError("product unique constraint violated")
		}
		return nil, fmt.Errorf("update product: %w", err)
	}

	res := dto.ProductToResponse(product)
	return &res, nil
}

func (u *productUseCase) Delete(ctx context.Context, req presentation.GetProductRequest) error {
	err := validation.Validate(&req)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return fmt.Errorf("parse id: %w", err)
	}

	err = u.repo.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return typex.NewNotFoundError("product")
		}
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}
