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
	var price int32 = 10

	mockedUseCase := mocks.UseCase{}
	server := NewProductAPIServer(&mockedUseCase)

	t.Run("list server with one product", func(t *testing.T) {
		products := []*entity.Product{
			{
				ID:           uuid.NewV4().String(),
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: price,
			},
		}
		grpcProducts := []*grpcv1.CreateProductResponse{
			{
				Id:           products[0].ID,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: price,
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
				PriceInCents: price,
			},
			{
				ID:           uuid.NewV4().String(),
				Title:        "cacilds vidis",
				Description:  "Todo mundo vê os porris que eu tomo",
				PriceInCents: price,
			},
		}
		grpcProducts := []*grpcv1.CreateProductResponse{
			{
				Id:           products[0].ID,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: price,
			},
			{
				Id:           products[1].ID,
				Title:        "cacilds vidis",
				Description:  "Todo mundo vê os porris que eu tomo",
				PriceInCents: price,
			},
		}
		request := grpcv1.ListProductsRequest{}
		ctx := context.TODO()

		mockedUseCase.On("List", mock.AnythingOfType("*context.valueCtx")).Return(products, nil)

		response, err := server.ListProducts(ctx, &request)

		assert.Nil(t, err)
		assert.Equal(t, grpcProducts, response.Products)
	})

	t.Run("create product", func(t *testing.T) {
		request := grpcv1.CreateProductRequest{
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: price,
		}
		mocked := entity.Product{
			ID:           uuid.NewV4().String(),
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: price,
		}
		expected := grpcv1.CreateProductResponse{
			Id:           mocked.ID,
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: price,
		}
		ctx := context.Background()

		mockedUseCase.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.Anything).Return(&mocked, nil)

		response, err := server.CreateProduct(ctx, &request)

		assert.Nil(t, err)
		assert.Equal(t, &expected, response)
	})
}
