package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

// Store represents a database interface
type Store struct {
	db *sqlx.DB
}

// NewStore constructs a database store
func NewStore(driver, dsn string) (*Store, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(-1)

	return &Store{db}, nil
}

func (st *Store) truncateTables(tables []string) error {
	_, err := st.db.Exec(
		fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", strings.Join(tables, ", ")),
	)

	return err
}

func (st *Store) BeginTx() (*sql.Tx, error) {
	return st.db.Begin()
}

func (st *Store) Migrate() error {
	return goose.Up(st.db.DB, "./database/migrations")
}
