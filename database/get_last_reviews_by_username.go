package database

import (
	"context"

	"github.com/ansel1/merry"
)

// GetLastReviewsByUsername returns N most recent reviews for a given user
func (st *Store) GetLastReviewsByUsername(ctx context.Context, username string, limit int) (*[]*ReviewRecord, error) {
	var r []*ReviewRecord
	err := st.db.SelectContext(ctx, &r, `
    SELECT * FROM reviews WHERE github_username = $1 ORDER BY created_at DESC LIMIT $2
  `, username, limit)

	if err != nil {
		return nil, merry.Here(err)
	}

	return &r, nil
}
