package service

import (
	"context"

	"github.com/mtavano/ghreviews"
	"github.com/mtavano/ghreviews/database"
	"github.com/sirupsen/logrus"
)

var _ ghreviews.ReviewService = &reviewService{}

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

func (r *reviewService) GetLastReviews() ([]*ghreviews.GhReview, error) {
	rr, err := r.store.GetLastReviews(context.Background(), 10)
	if err != nil {
		return nil, err
	}

	reviews := make([]*ghreviews.GhReview, len(*rr))
	for i := range *rr {
		reviews[i] = (*rr)[i].ToGhReview()
	}

	return reviews, nil
}

func (r *reviewService) GetLastReviewsByUsername(githubUsername string) ([]*ghreviews.GhReview, error) {
	rr, err := r.store.GetLastReviewsByUsername(context.Background(), githubUsername, 10)
	if err != nil {
		return nil, err
	}

	reviews := make([]*ghreviews.GhReview, len(*rr))
	for i := range *rr {
		reviews[i] = (*rr)[i].ToGhReview()
	}

	return reviews, nil
}

func (r *reviewService) CountReviews() (int, error) {
	count, err := r.store.CountReviews(context.Background())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *reviewService) CountReviewsByUsername(githubUsername string) (int, error) {
	count, err := r.store.CountReviewsByUsername(context.Background(), githubUsername)
	if err != nil {
		return 0, err
	}

	return count, nil
}
