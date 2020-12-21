package service

import (
	"context"

	"github.com/mtavano/ghreviews"
	"github.com/mtavano/ghreviews/database"
	"github.com/sirupsen/logrus"
)

type reviewService struct {
	store  *database.Store
	logger *logrus.Logger
}

func NewReviewService(logger *logrus.Logger, store *database.Store) *reviewService {
	return &reviewService{
		store,
		logger,
	}
}

func (r *reviewService) CreateReview(githubUsername, githubAvatarUrl, content string) (*ghreviews.GhReview, error) {
	review, err := r.store.CreateReview(context.Background(), githubUsername, githubAvatarUrl, content, nil)
	if err != nil {
		return nil, err
	}

	return review.ToGhReview(), nil
}
