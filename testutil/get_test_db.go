package testutil

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// NewTestDatabase returns a new test db with all migrations applied
func NewTestDatabase(t *testing.T) *sqlx.DB {
	t.Helper()

	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	db.SetConnMaxLifetime(-1)

	return db
}

