package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upProvisioning, downProvisioning)
}

func upProvisioning(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
    CREATE EXTENSION pgcrypto;
  `)
	if err != nil {
		return err
	}
	return nil
}

func downProvisioning(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
    DROP EXTENSION pgcrypto;
  `)
	if err != nil {
		return err
	}
	return nil
}
