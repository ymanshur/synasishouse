package usecase

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
	"github.com/ymanshur/synasishouse/inventory/internal/consts"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	mockrepo "github.com/ymanshur/synasishouse/inventory/internal/repo/mock"
	"github.com/ymanshur/synasishouse/inventory/internal/typex"
)

func TestProduct_Create(t *testing.T) {
	product := randomProduct(0)

	testCases := []struct {
		name          string
		req           presentation.CreateProductRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.ProductResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.CreateProductRequest{
				Code:  product.Code,
				Total: product.Total,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CreateProduct(gomock.Any(), db.CreateProductParams{
						Code:  product.Code,
						Total: product.Total,
					}).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotZero(t, res.ID)
				require.Equal(t, product.Code, res.Code)
				require.Equal(t, product.Total, res.Total)
				require.Zero(t, product.Reserved)
			},
		},
		{
			name: "DuplicateCode",
			req: presentation.CreateProductRequest{
				Code:  product.Code,
				Total: product.Total,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CreateProduct(gomock.Any(), db.CreateProductParams{
						Code:  product.Code,
						Total: product.Total,
					}).
					Times(1).
					Return(db.Product{}, &pgconn.PgError{Code: consts.UniqueViolation})
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
				var unprocessableEnityErr typex.UnProcessableEnity
				require.ErrorAs(t, err, &unprocessableEnityErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewProduct(repo)
			res, err := u.Create(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestProduct_Get(t *testing.T) {
	product := randomProduct(10)

	testCases := []struct {
		name          string
		req           presentation.GetProductRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.ProductResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.GetProductRequest{
				ID: product.ID.String(),
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					GetProduct(gomock.Any(), product.ID).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotZero(t, res.ID)
				require.Equal(t, product.Code, res.Code)
				require.Equal(t, product.Total, res.Total)
				require.Equal(t, product.Reserved, int32(10))
			},
		},
		{
			name: "NotFound",
			req: presentation.GetProductRequest{
				ID: product.ID.String(),
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					GetProduct(gomock.Any(), product.ID).
					Times(1).
					Return(db.Product{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
		{
			name: "InvalidID",
			req: presentation.GetProductRequest{
				ID: "123",
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					GetProduct(gomock.Any(), product.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewProduct(repo)
			res, err := u.Get(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestProduct_Update(t *testing.T) {
	product := randomProduct(0)

	testCases := []struct {
		name          string
		req           presentation.UpdateProductRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.ProductResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.UpdateProductRequest{
				ID:   product.ID.String(),
				Code: product.Code,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					UpdateProduct(gomock.Any(), db.UpdateProductParams{
						Code: product.Code,
						ID:   product.ID,
					}).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotZero(t, res.ID)
				require.Equal(t, product.Code, res.Code)
			},
		},
		{
			name: "DuplicateCode",
			req: presentation.UpdateProductRequest{
				ID:   product.ID.String(),
				Code: product.Code,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					UpdateProduct(gomock.Any(), db.UpdateProductParams{
						Code: product.Code,
						ID:   product.ID,
					}).
					Times(1).
					Return(db.Product{}, &pgconn.PgError{Code: consts.UniqueViolation})
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
				var unprocessableEnityErr typex.UnProcessableEnity
				require.ErrorAs(t, err, &unprocessableEnityErr)
			},
		},
		{
			name: "NotFound",
			req: presentation.UpdateProductRequest{
				ID:   product.ID.String(),
				Code: product.Code,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					UpdateProduct(gomock.Any(), db.UpdateProductParams{
						Code: product.Code,
						ID:   product.ID,
					}).
					Times(1).
					Return(db.Product{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
		{
			name: "InvalidID",
			req: presentation.UpdateProductRequest{
				Code: product.Code,
				ID:   "123",
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewProduct(repo)
			res, err := u.Update(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestProduct_Delete(t *testing.T) {
	product := randomProduct(10)

	testCases := []struct {
		name          string
		req           presentation.GetProductRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.ProductResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.GetProductRequest{
				ID: product.ID.String(),
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					DeleteProduct(gomock.Any(), product.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			req: presentation.GetProductRequest{
				ID: product.ID.String(),
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					DeleteProduct(gomock.Any(), product.ID).
					Times(1).
					Return(pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
				var notFoundErr typex.NotFound
				require.ErrorAs(t, err, &notFoundErr)
			},
		},
		{
			name: "InvalidID",
			req: presentation.GetProductRequest{
				ID: "123",
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					DeleteProduct(gomock.Any(), product.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockrepo.NewMockRepo(c)

			tc.buildStubs(repo)

			u := NewProduct(repo)
			err := u.Delete(context.Background(), tc.req)
			tc.checkResponse(t, nil, err)
		})
	}
}

func randomProduct(reserved int32) db.Product {
	return db.Product{
		ID:       uuid.New(),
		Code:     randomProductCode(),
		Total:    randomInt(1, 100),
		Reserved: reserved,
	}
}

func randomProductCode() string {
	no := randomInt(1, 200)
	return fmt.Sprintf("P%d", no)
}

func randomAmount() int32 {
	return randomInt(1000, 100000)
}

func randomInt(min, max int32) int32 {
	return min + rand.Int32N(max-min+1)
}
