package database

import (
	"context"

	"github.com/ansel1/merry"
)

// GetLastReviews returns N most recent reviews
func (st *Store) GetLastReviews(ctx context.Context, limit int) (*[]*ReviewRecord, error) {
	var r []*ReviewRecord
	err := st.db.SelectContext(ctx, &r, `
    SELECT * FROM reviews ORDER BY created_at DESC LIMIT $1
  `, limit)

	if err != nil {
		return nil, merry.Here(err)
	}

	return &r, nil
}
