package database

import (
	"context"

	"github.com/ansel1/merry"
)

// CountReviewsByUsername returns total amount of reviews for a given user
func (st *Store) CountReviewsByUsername(ctx context.Context, username string) (int, error) {
	var count int
	err := st.db.GetContext(ctx, &count, `
    SELECT count(*) FROM reviews WHERE github_username = $1
  `, username)

	if err != nil {
		return count, merry.Here(err)
	}

	return count, nil
}
