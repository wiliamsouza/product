package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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
			"struct=postgres.ProductDataStore, method=List, error=select_fail",
		)
		return products, wrappedErr
	}

	return products, nil
}

// Create product.
func (p *ProductDataStore) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	ctx, span := trace.StartSpan(ctx, "postgres.ProductDataStore.Create")
	defer span.End()
	queryStatement := `
		INSERT
			INTO
			product (id,
			title,
			description,
			price_in_cents)
		VALUES ($1,
		$2,
		$3,
		$4) RETURNING id,
		title,
		description,
		price_in_cents
	`
	var result entity.Product
	err := p.DB.QueryRowxContext(ctx, queryStatement, uuid.NewV4(), product.Title, product.Description, product.PriceInCents).StructScan(&result)
	if err != nil {
		wrappedErr := errors.Wrapf(
			err,
			"struct=postgres.ProductDataStore, method=Create, error=insert_fail",
		)
		return &result, wrappedErr
	}

	return &result, nil
}

//NewProductDataStore create a product data store instance.
func NewProductDataStore(db *sqlx.DB) *ProductDataStore {
	return &ProductDataStore{DB: db}
}
