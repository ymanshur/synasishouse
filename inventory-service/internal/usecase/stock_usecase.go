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
	Check(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error)
	Reserve(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error)
	Release(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error)
}

type stockUseCase struct {
	repo repo.Repo
}

func NewStock(repo repo.Repo) Stocker {
	return &stockUseCase{repo: repo}
}

func (u *stockUseCase) Check(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	params := repo.CreateStockParams{Stocks: []repo.StockParams{}}
	for _, stock := range req.Stocks {
		params.Stocks = append(params.Stocks, repo.StockParams{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}

	isAvailable, err := u.repo.CheckStock(ctx, params)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, typex.NewNotFoundError("stock")
		}
		return nil, fmt.Errorf("check stock: %w", err)
	}

	res := &presentation.StockResponse{IsAvailable: isAvailable}
	return res, nil
}

func (u *stockUseCase) Reserve(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	params := repo.CreateStockParams{Stocks: []repo.StockParams{}}
	for _, stock := range req.Stocks {
		params.Stocks = append(params.Stocks, repo.StockParams{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}

	isAvailable, err := u.repo.ReserveStock(ctx, params)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, typex.NewNotFoundError("stock")
		}
		return nil, fmt.Errorf("reserve stock: %w", err)
	}

	res := &presentation.StockResponse{IsAvailable: isAvailable}
	return res, nil
}

func (u *stockUseCase) Release(ctx context.Context, req presentation.CreateStockRequest) (*presentation.StockResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	params := repo.CreateStockParams{Stocks: []repo.StockParams{}}
	for _, stock := range req.Stocks {
		params.Stocks = append(params.Stocks, repo.StockParams{
			ProductCode: stock.ProductCode,
			Amount:      stock.Amount,
		})
	}

	isAvailable, err := u.repo.ReleaseStock(ctx, params)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, typex.NewNotFoundError("stock")
		}
		return nil, fmt.Errorf("release stock: %w", err)
	}

	res := &presentation.StockResponse{IsAvailable: isAvailable}
	return res, nil
}
