package usecase

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"wiliam.dev/product"
	"wiliam.dev/product/entity"
)

// Ensure ProductUseCase implements product.UseCase interface.
var _ product.UseCase = &ProductUseCase{}

//ProductUseCase implements product.UseCase interface.
type ProductUseCase struct {
	DataStore product.DataStore
}

// List products.
func (u *ProductUseCase) List(ctx context.Context) ([]*entity.Product, error) {
	ctx, span := trace.StartSpan(ctx, "usecase.ProductUseCase.List")
	defer span.End()
	products, err := u.DataStore.List(ctx)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			wrappedErr := errors.Wrapf(
				errors.New("resource not found"),
				"struct=usecase.ProductUseCase, method=List, error=repository_list_not_found",
			)
			return nil, wrappedErr
		}
		wrappedErr := errors.Wrapf(
			err,
			"struct=usecase.ProductUseCase, method=List, error=repository_list_failed",
		)
		return nil, wrappedErr
	}
	return products, nil
}

//NewProductUseCase create a product use case instance.
func NewProductUseCase(dataStore product.DataStore) *ProductUseCase {
	return &ProductUseCase{DataStore: dataStore}
}
