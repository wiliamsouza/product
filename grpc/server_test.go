package grpc

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wiliam.dev/product/entity"
	grpcv1 "wiliam.dev/product/grpc/v1beta1"
	"wiliam.dev/product/mocks"
)

func TestGRPCServer(t *testing.T) {

	mockedUseCase := mocks.UseCase{}
	server := NewProductAPIServer(&mockedUseCase)

	t.Run("list server with one product", func(t *testing.T) {
		products := []*entity.Product{
			{
				ID:           uuid.NewV4().String(),
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: 10,
			},
		}
		grpcProducts := []*grpcv1.Product{
			{
				Id:           products[0].ID,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: 10,
			},
		}
		request := grpcv1.ListProductsRequest{}
		ctx := context.Background()

		mockedUseCase.On("List", mock.AnythingOfType("*context.valueCtx")).Return(products, nil)

		response, err := server.ListProducts(ctx, &request)

		assert.Nil(t, err)
		assert.Equal(t, grpcProducts, response.Products)
	})

	t.Run("list server with multiple products", func(t *testing.T) {
		mockedUseCase.Mock = mock.Mock{} // Reset mock
		products := []*entity.Product{
			{
				ID:           uuid.NewV4().String(),
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: 10,
			},
			{
				ID:           uuid.NewV4().String(),
				Title:        "cacilds vidis",
				Description:  "Todo mundo vê os porris que eu tomo",
				PriceInCents: 15,
			},
		}
		grpcProducts := []*grpcv1.Product{
			{
				Id:           products[0].ID,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: 10,
			},
			{
				Id:           products[1].ID,
				Title:        "cacilds vidis",
				Description:  "Todo mundo vê os porris que eu tomo",
				PriceInCents: 15,
			},
		}
		request := grpcv1.ListProductsRequest{}
		ctx := context.TODO()

		mockedUseCase.On("List", mock.AnythingOfType("*context.valueCtx")).Return(products, nil)

		response, err := server.ListProducts(ctx, &request)

		assert.Nil(t, err)
		assert.Equal(t, grpcProducts, response.Products)
	})
}
