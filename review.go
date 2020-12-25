package ghreviews

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	// Badges

	// BadgeMaster is the master badge for a review
	BadgeMaster = "master"
	// BadgeCool is the cool badge for a review
	BadgeCool = "cool"
	// BadgeTalented is the talented badge for a review
	BadgeTalented = "talented"
)

type ReviewService interface {
	CreateReview(githubUsername, githubAvatarUrl, content string) (*GhReview, error)
	GetLastReviews() ([]*GhReview, error)
	GetLastReviewsByUsername(username string) ([]*GhReview, error)
	CountReviews() (int, error)
	CountReviewsByUsername(username string) (int, error)
}

// GhReview represents a github review app model
type GhReview struct {
	ID              string  `json:"id"`
	GithubUsername  string  `json:"githubUsername"`
	GithubAvatarURL string  `json:"githubAvatarUrl"`
	Content         string  `json:"content"`
	Badge           *string `json:"badge"`
	CreatedAt       int64   `json:"createdAt"`
}

type CreateReviewInput struct {
	GithubUsername  string `json:"githubUsername" validate:"required"`
	GithubAvatarURL string `json:"githubAvatarUrl"`
	Content         string `json:"content" validate:"required"`
}

type ServiceError struct {
	Code        uint
	Description string
	Err         error
}

func (e ServiceError) Error() string {
	return e.Description
}

func (e ServiceError) Unwrap() error {
	return e.Err
}

func NewNotFoundErr(ctx context.Context, err error) *gqlerror.Error {
	return &gqlerror.Error{Message: err.Error(), Path: graphql.GetPath(ctx), Extensions: map[string]interface{}{
		"code":        404,
		"description": "Resource not found",
	}}
}
