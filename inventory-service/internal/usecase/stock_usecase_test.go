package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	"github.com/ymanshur/synasishouse/inventory/internal/repo"
	mockrepo "github.com/ymanshur/synasishouse/inventory/internal/repo/mock"
	"github.com/ymanshur/synasishouse/inventory/internal/typex"
)

func TestStock_Check(t *testing.T) {
	product := randomProduct(0)
	amount := product.Total - randomInt(0, 999)

	testCases := []struct {
		name          string
		req           presentation.CreateStockRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.StockResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CheckStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.True(t, res.IsAvailable)
			},
		},
		{
			name: "UnAvailableStock",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CheckStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.False(t, res.IsAvailable)
			},
		},
		{
			name: "NotFound",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CheckStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewStock(repo)
			res, err := u.Check(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestStock_Reserve(t *testing.T) {
	product := randomProduct(0)
	amount := product.Total - randomInt(0, 999)

	testCases := []struct {
		name          string
		req           presentation.CreateStockRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.StockResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReserveStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.True(t, res.IsAvailable)
			},
		},
		{
			name: "UnAvailableReserve",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReserveStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.False(t, res.IsAvailable)
			},
		},
		{
			name: "NotFound",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReserveStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewStock(repo)
			res, err := u.Reserve(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestStock_Release(t *testing.T) {
	product := randomProduct(0)
	amount := product.Total - randomInt(0, 999)

	testCases := []struct {
		name          string
		req           presentation.CreateStockRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.StockResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReleaseStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.True(t, res.IsAvailable)
			},
		},
		{
			name: "UnAvailableRelease",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReleaseStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.False(t, res.IsAvailable)
			},
		},
		{
			name: "NotFound",
			req: presentation.CreateStockRequest{
				Stocks: []presentation.StockRequest{
					{
						ProductCode: product.Code,
						Amount:      amount,
					},
				},
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					ReleaseStock(gomock.Any(), repo.CreateStockParams{
						Stocks: []repo.StockParams{
							{
								ProductCode: product.Code,
								Amount:      amount,
							},
						},
					}).
					Times(1).
					Return(false, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.StockResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewStock(repo)
			res, err := u.Release(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
