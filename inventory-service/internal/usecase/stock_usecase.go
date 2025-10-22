package usecase

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ymanshur/synasishouse/inventory/internal/consts"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/inventory/internal/repo"
	"github.com/ymanshur/synasishouse/inventory/internal/typex"
)

type Stocker interface {
	Check(ctx context.Context, req presentation.StockRequest) error
	Reserve(ctx context.Context, req presentation.StockRequest) error
	Release(ctx context.Context, req presentation.StockRequest) error
}

type stockUseCase struct {
	repo repo.Repo
}

func NewStock(repo repo.Repo) Stocker {
	return &stockUseCase{repo: repo}
}

func (u *stockUseCase) Check(ctx context.Context, req presentation.StockRequest) error {
	err := validation.Validate(&req)
	if err != nil {
		return err
	}

	err = u.repo.CheckStock(ctx, req.Code, req.Amount)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return typex.NewNotFoundError("stock")
		}
		return fmt.Errorf("check stock: %w", err)
	}

	return nil
}

func (u *stockUseCase) Reserve(ctx context.Context, req presentation.StockRequest) error {
	err := validation.Validate(&req)
	if err != nil {
		return err
	}

	err = u.repo.ReserveStock(ctx, req.Code, req.Amount)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return typex.NewNotFoundError("stock")
		}
		return fmt.Errorf("reserve stock: %w", err)
	}

	return nil
}

func (u *stockUseCase) Release(ctx context.Context, req presentation.StockRequest) error {
	err := validation.Validate(&req)
	if err != nil {
		return err
	}

	err = u.repo.ReleaseStock(ctx, req.Code, req.Amount)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return typex.NewNotFoundError("stock")
		}
		return fmt.Errorf("release stock: %w", err)
	}

	return nil
}
