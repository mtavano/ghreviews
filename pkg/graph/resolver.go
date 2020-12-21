package graph

import (
	"context"
	"sync"

	"github.com/mtavano/ghreviews"
	"github.com/sirupsen/logrus"
)

type Resolver struct {
	reviewService ghreviews.ReviewService
	logger        *logrus.Logger

	mu           sync.RWMutex
	reviewersHub map[chan []*ghreviews.GhReview]bool
}

func NewResolver(logger *logrus.Logger, reviewService ghreviews.ReviewService) *Resolver {
	return &Resolver{
		reviewService: reviewService,
		logger:        logger,
		reviewersHub:  make(map[chan []*ghreviews.GhReview]bool),
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

	msg := []*ghreviews.GhReview{review}
	go func() {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for channel := range r.reviewersHub {
			channel <- msg
		}
	}()

	return review, nil
}

func (r *queryResolver) GetReview(ctx context.Context, id string) (*ghreviews.GhReview, error) {
	r.logger.Debugln("Not implemented")
	return nil, nil
}

func (r *subscriptionResolver) Feed(ctx context.Context) (<-chan []*ghreviews.GhReview, error) {
	r.mu.RLock()
	r.logger.Debugln("Connected: ", len(r.reviewersHub), " users")
	r.mu.RUnlock()
	rr, err := r.reviewService.GetLastReviews()
	if err != nil {
		return nil, err
	}

	// Add channel to hub
	cr := make(chan []*ghreviews.GhReview, 1)
	r.mu.Lock()
	r.reviewersHub[cr] = true
	r.mu.Unlock()
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.reviewersHub, cr)
		r.mu.UnLock()
	}()

	cr <- rr

	return cr, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Subscription returns MutationResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
