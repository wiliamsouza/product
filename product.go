package product

import (
	"context"

	"wiliam.dev/product/entity"
)

// DataStore store interface for product entities.
type DataStore interface {
	List(ctx context.Context) ([]*entity.Product, error)
}

// UseCase domain interface for product bussines logic.
type UseCase interface {
	List(ctx context.Context) ([]*entity.Product, error)
}
