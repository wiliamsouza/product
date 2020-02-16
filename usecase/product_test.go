package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wiliam.dev/product/entity"
	"wiliam.dev/product/mocks"
)

func TestProductUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	useCase := NewProductUseCase(&mockedDataStore)

	t.Run("Test List", func(t *testing.T) {
		ctx := context.TODO()
		products := []*entity.Product{}
		mockedDataStore.On("List", mock.AnythingOfType("*context.valueCtx")).Return([]*entity.Product{}, nil)
		productsFromDataStore, err := useCase.List(ctx)
		assert.Nil(t, err)
		assert.Equal(t, products, productsFromDataStore)
	})
}

func TestNewProductUseCase(t *testing.T) {
	mockedDataStore := mocks.DataStore{}
	product := NewProductUseCase(&mockedDataStore)
	assert.IsType(t, &ProductUseCase{}, product)
	assert.NotNil(t, product.DataStore)
}
