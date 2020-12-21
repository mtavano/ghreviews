package database

import (
	"context"

	"github.com/ansel1/merry"
)

// CreateReview inserts a review on the database
func (st *Store) CreateReview(ctx context.Context, githubUsername, githubAvatarUrl, content string, badge *string) (*ReviewRecord, error) {
	var r ReviewRecord
	err := st.db.GetContext(ctx, &r, `
    INSERT INTO reviews (github_username, github_avatar_url, content, badge)
      VALUES ($1, $2, $3, $4)
      RETURNING *;
  `, githubUsername, githubAvatarUrl, content, badge)

	if err != nil {
		return nil, merry.Here(err)
	}

	return &r, nil
}
