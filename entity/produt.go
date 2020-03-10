package entity

// Product entity.
type Product struct {
	ID           string
	Title        string `db:"title"`
	Description  string `db:"description"`
	PriceInCents int32  `db:"price_in_cents"`
	Discount     Discount
}

// Discount entity.
type Discount struct {
	Pct          float32
	ValueInCents int32
}
