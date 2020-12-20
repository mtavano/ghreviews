package ghreview

import "time"

var (
	// Badges

	// BadgeMaster is the master badge for a review
	BadgeMaster = "master"
	// BadgeCool is the cool badge for a review
	BadgeCool = "cool"
	// BadgeTalented is the talented badge for a review
	BadgeTalented = "talented"
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
func (rr *ReviewRecord) ToGhReview() *GhReview {
	return &GhReview{
		GithubUsername:  rr.GithubUsername,
		GithubAvatarURL: rr.GithubAvatarURL,
		Content:         rr.Content,
		Badge:           rr.Badge,
	}
}

// GhReview represents a github review app model
type GhReview struct {
	GithubUsername  string  `json:"github_username" validate:"required"`
	GithubAvatarURL string  `json:"github_avatar_url"`
	Content         string  `json:"content" validate:"required"`
	Badge           *string `json:"badge" validate:"omitempty,eq=master|cool|talented"`
}
