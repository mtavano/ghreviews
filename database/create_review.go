package database

import (
	"context"

	"github.com/ansel1/merry"
	ghreview "github.com/mtavano/ghreviews"
)

// CreateReview inserts a review on the database
func (st *Store) CreateReview(ctx context.Context, review *ghreview.GhReview) (*ghreview.ReviewRecord, error) {
	var r ghreview.ReviewRecord
	err := st.db.GetContext(ctx, &r, `
    INSERT INTO reviews (github_username, content, badge)
      VALUES ($1, $2, $3)
      RETURNING *;
  `, review.GithubUsername, review.Content, review.Badge)

	if err != nil {
		return nil, merry.Here(err)
	}

	return &r, nil
}
