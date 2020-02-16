package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductDataStore(t *testing.T) {

	queryStatement := `
		SELECT
			id,
			title,
			description,
			price_in_cents
		FROM
			product
	`

	t.Run("list with one row in database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "price_in_cents"}).
			AddRow(
				uuid.NewV4(),
				"cacilds vidis litro abertis",
				"Todo mundo vê os porris que eu tomo",
				50,
			)

		mock.ExpectQuery(queryStatement).WillReturnRows(rows)

		ctx := context.TODO()

		db := sqlx.NewDb(mockDB, "sqlmock")
		defer db.Close()
		r := NewProductDataStore(db)
		products, err := r.List(ctx)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(products))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("list with multiple rows in database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "price_in_cents"}).
			AddRow(
				uuid.NewV4(),
				"cacilds vidis litro abertis",
				"Todo mundo vê os porris que eu tomo",
				50,
			).
			AddRow(
				uuid.NewV4(),
				"Mussum Ipsum",
				"mas ninguém vê os tombis que eu levo",
				100,
			)

		mock.ExpectQuery(queryStatement).WillReturnRows(rows)

		ctx := context.TODO()

		db := sqlx.NewDb(mockDB, "sqlmock")
		defer db.Close()
		r := NewProductDataStore(db)
		products, err := r.List(ctx)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(products))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
