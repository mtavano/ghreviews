package database

import (
	"context"

	"github.com/ansel1/merry"
)

// CountReviews returns total amount of reviews for a given user
func (st *Store) CountReviews(ctx context.Context) (int, error) {
	var count int
	err := st.db.GetContext(ctx, &count, `
    SELECT count(*) FROM reviews
  `)

	if err != nil {
		return count, merry.Here(err)
	}

	return count, nil
}
