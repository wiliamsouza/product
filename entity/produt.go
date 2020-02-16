package entity

// Product entity.
type Product struct {
	ID           string
	Title        string
	Description  string
	PriceInCents int32 `db:"price_in_cents"`
}
