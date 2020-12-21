package graph

import (
	"context"

	"github.com/mtavano/ghreviews"
	"github.com/sirupsen/logrus"
)

type Resolver struct {
	reviewService ghreviews.ReviewService
	logger        *logrus.Logger
}

func NewResolver(logger *logrus.Logger, reviewService ghreviews.ReviewService) *Resolver {
	return &Resolver{
		reviewService,
		logger,
	}
}

func (r *mutationResolver) CreateReview(ctx context.Context, reviewInput ghreviews.CreateReviewInput) (*ghreviews.GhReview, error) {
	review, err := r.reviewService.CreateReview(
		reviewInput.GithubUsername,
		reviewInput.GithubAvatarURL,
		reviewInput.Content,
	)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *queryResolver) GetReview(ctx context.Context, id string) (*ghreviews.GhReview, error) {
	r.logger.Debugln("Not implemented")
	return nil, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
