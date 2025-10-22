package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
	mockrepo "github.com/ymanshur/synasishouse/inventory/internal/repo/mock"
)

func TestProduct_Create(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	testCases := []struct {
		name          string
		req           presentation.CreateProductRequest
		buildStubs    func(mr *mockrepo.MockRepo)
		checkResponse func(t *testing.T, res *presentation.ProductResponse, err error)
	}{
		{
			name: "Success",
			req: presentation.CreateProductRequest{
				Code:  "C001",
				Total: 10,
			},
			buildStubs: func(mr *mockrepo.MockRepo) {
				mr.EXPECT().
					CreateProduct(gomock.Any(), db.CreateProductParams{
						Code:  "C001",
						Total: 10,
					}).
					Times(1).
					Return(db.Product{}, nil)
			},
			checkResponse: func(t *testing.T, res *presentation.ProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
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
