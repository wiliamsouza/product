package usecase

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wiliam.dev/product/entity"
	"wiliam.dev/product/mocks"
)

func TestProductUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	useCase := NewProductUseCase(&mockedDataStore)

	t.Run("Test list", func(t *testing.T) {
		ctx := context.TODO()
		products := []*entity.Product{}
		mockedDataStore.On("List", mock.AnythingOfType("*context.valueCtx")).Return([]*entity.Product{}, nil)
		productsFromDataStore, err := useCase.List(ctx)
		assert.Nil(t, err)
		assert.Equal(t, products, productsFromDataStore)
	})

	t.Run("Test create", func(t *testing.T) {
		ctx := context.TODO()
		product := entity.Product{
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: 10,
		}
		mocked := entity.Product{
			ID:           uuid.NewV4().String(),
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: 10,
		}
		expected := entity.Product{
			ID:           mocked.ID,
			Title:        "Mussum Ipsum",
			Description:  "cacilds vidis litro abertis",
			PriceInCents: 10,
		}
		mockedDataStore.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.Anything).Return(&mocked, nil)
		productFromDataStore, err := useCase.Create(ctx, &product)
		assert.Nil(t, err)
		assert.Equal(t, &expected, productFromDataStore)
	})
}

func TestNewProductUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	product := NewProductUseCase(&mockedDataStore)
	assert.IsType(t, &ProductUseCase{}, product)
	assert.NotNil(t, product.DataStore)
}
