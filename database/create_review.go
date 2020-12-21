package database

import (
	"context"
	"time"

	"github.com/ansel1/merry"
	"github.com/mtavano/ghreviews"
)

// ReviewRecord represents a database record
type ReviewRecord struct {
	ID              string    `db:"id"`
	GithubUsername  string    `db:"github_username"`
	GithubAvatarURL string    `db:"github_avatar_url"`
	Content         string    `db:"content"`
	Badge           *string   `db:"badge"`
	CreatedAt       time.Time `db:"created_at"`
}

// ToGhReview returns a ToGhReview pointer
func (rr *ReviewRecord) ToGhReview() *ghreviews.GhReview {
	return &ghreviews.GhReview{
		ID:              rr.ID,
		GithubUsername:  rr.GithubUsername,
		GithubAvatarURL: rr.GithubAvatarURL,
		Content:         rr.Content,
		Badge:           rr.Badge,
		CreatedAt:       toMilliseconds(rr.CreatedAt),
	}
}

func toMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

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
