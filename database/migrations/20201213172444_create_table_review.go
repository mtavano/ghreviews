package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableReview, downCreateTableReview)
}

func upCreateTableReview(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
    CREATE TYPE  badge AS ENUM ('master', 'cool', 'talented');

    CREATE TABLE reviews (
      id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      github_username       TEXT NOT NULL,
      github_avatar_url     TEXT NOT NULL,
      content               TEXT NOT NULL,
      badge                 badge,
      created_at            timestamptz DEFAULT NOW()
    );
  `)
	return err
}

func downCreateTableReview(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
    DROP TABLE reviews;
    DROP TABLE badge;
  `)
	return err
}
