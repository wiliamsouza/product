package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"wiliam.dev/product"
	"wiliam.dev/product/entity"
)

// Ensure ProductDataStore implements product.DataStore interface.
var _ product.DataStore = &ProductDataStore{}

//ProductDataStore implements product.DataStore interface.
type ProductDataStore struct {
	DB *sqlx.DB
}

// List products from database.
func (p *ProductDataStore) List(ctx context.Context) ([]*entity.Product, error) {
	ctx, span := trace.StartSpan(ctx, "postgres.ProductDataStore.List")
	defer span.End()
	var products = []*entity.Product{}
	queryStatement := `
		SELECT
			id,
			title,
			description,
			price_in_cents
		FROM
			product
	`
	err := p.DB.SelectContext(ctx, &products, queryStatement)
	if err != nil {
		wrappedErr := errors.Wrapf(
			err,
			"struct=postgres.ProductProduct, method=List, error=select_fail",
		)
		return products, wrappedErr
	}

	return products, nil
}

//NewProductDataStore create a product data store instance.
func NewProductDataStore(db *sqlx.DB) *ProductDataStore {
	return &ProductDataStore{DB: db}
}
