package usecase

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	"wiliam.dev/product/entity"
	promotionv1 "wiliam.dev/product/grpc/client/promotion/v1alpha1"
	"wiliam.dev/product/mocks"
)

func TestPromtionUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	productUseCase := NewProductUseCase(&mockedDataStore)

	mockedPromotionClient := mocks.PromotionAPIClient{}
	useCase := NewPromotionUseCase(productUseCase, &mockedPromotionClient)

	var (
		price         int32   = 10000
		percentage    float32 = 5
		discountValue int32   = 500
	)

	t.Run("Test list empty", func(t *testing.T) {
		ctx := context.TODO()
		expected := []*entity.Product{}
		mockedDataStore.On("List", mock.AnythingOfType("*context.valueCtx")).Return([]*entity.Product{}, nil)
		productsFromDataStore, err := useCase.List(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expected, productsFromDataStore)
	})

	t.Run("Test list product with discount", func(t *testing.T) {
		mockedDataStore.Mock = mock.Mock{} // Reset mock
		ctx := context.TODO()
		m := make(map[string]string)
		m["x-user-id"] = uuid.NewV4().String()
		md := metadata.New(m)
		ctxUser := metadata.NewIncomingContext(ctx, md)
		id := uuid.NewV4().String()
		mocked := []*entity.Product{
			{
				ID:           id,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: price,
			},
		}
		expected := []*entity.Product{
			{
				ID:           id,
				Title:        "Mussum Ipsum",
				Description:  "cacilds vidis litro abertis",
				PriceInCents: price,
				Discount: entity.Discount{
					Pct:          percentage,
					ValueInCents: discountValue,
				},
			},
		}
		mockedResponse := promotionv1.RetrievePromotionResponse{
			Discounts: []*promotionv1.Discount{
				{
					Pct: percentage,
				},
			},
		}
		mockedDataStore.On("List", mock.AnythingOfType("*context.valueCtx")).Return(mocked, nil)
		mockedPromotionClient.On(
			"RetrievePromotion",
			mock.AnythingOfType("*context.valueCtx"),
			mock.Anything).Return(&mockedResponse, nil)
		productsFromDataStore, err := useCase.List(ctxUser)
		assert.Nil(t, err)
		assert.Equal(t, expected, productsFromDataStore)
	})
}

func TestNewPromotionUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	productUseCase := NewProductUseCase(&mockedDataStore)

	mockedPromotionClient := mocks.PromotionAPIClient{}
	promotion := NewPromotionUseCase(productUseCase, &mockedPromotionClient)

	assert.IsType(t, &ProductUseCase{}, promotion.Product)
	assert.NotNil(t, promotion.Promotion)
}
