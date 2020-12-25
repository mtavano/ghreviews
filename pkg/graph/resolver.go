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

	publicHubMU sync.RWMutex
	publicHub   map[chan []*ghreviews.GhReview]bool

	privateHubMU sync.RWMutex
	privateHub   map[string][](chan []*ghreviews.GhReview)
}

func NewResolver(logger *logrus.Logger, reviewService ghreviews.ReviewService) *Resolver {
	return &Resolver{
		reviewService: reviewService,
		logger:        logger,
		publicHub:     make(map[chan []*ghreviews.GhReview]bool),
		privateHub:    make(map[string][]chan []*ghreviews.GhReview),
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
		r.publicHubMU.RLock()
		defer r.publicHubMU.RUnlock()
		for channel := range r.publicHub {
			channel <- msg
		}
	}()
	go func() {
		r.privateHubMU.RLock()
		defer r.privateHubMU.RUnlock()
		for _, channel := range r.privateHub[reviewInput.GithubUsername] {
			channel <- msg
		}
	}()

	return review, nil
}

func (r *queryResolver) GetReview(ctx context.Context, id string) (*ghreviews.GhReview, error) {
	r.logger.Debugln("Not implemented")
	return nil, nil
}

func (r *subscriptionResolver) FeedByUsername(ctx context.Context, username string) (<-chan []*ghreviews.GhReview, error) {
	rr, err := r.reviewService.GetLastReviewsByUsername(username)
	if err != nil {
		return nil, err
	}

	// Add channel to private hub (by username)
	cr := make(chan []*ghreviews.GhReview, 1)
	r.privateHubMU.Lock()
	// if r.privateHub[username] == nil {
	// 	r.privateHub[username] = []chan []*ghreviews.GhReview{}
	// }
	r.privateHub[username] = append(r.privateHub[username], cr)
	r.privateHubMU.Unlock()
	go func() {
		<-ctx.Done()
		r.privateHubMU.Lock()
		newList := []chan []*ghreviews.GhReview{}
		for _, c := range r.privateHub[username] {
			if c == cr {
				continue
			}
			newList = append(newList, c)
		}
		r.privateHub[username] = newList
		r.privateHubMU.Unlock()
		r.logger.Debugf("Event: 'disconnected from %q'; %d users", username, len(r.privateHub[username]))
	}()

	cr <- rr
	r.publicHubMU.RLock()
	r.logger.Debugf("Event: 'connected to %q'; %d users", username, len(r.privateHub[username]))
	r.publicHubMU.RUnlock()

	return cr, nil
}

func (r *subscriptionResolver) Feed(ctx context.Context) (<-chan []*ghreviews.GhReview, error) {
	rr, err := r.reviewService.GetLastReviews()
	if err != nil {
		return nil, err
	}

	// Add channel to hub
	cr := make(chan []*ghreviews.GhReview, 1)
	r.publicHubMU.Lock()
	r.publicHub[cr] = true
	r.publicHubMU.Unlock()
	go func() {
		<-ctx.Done()
		r.publicHubMU.Lock()
		delete(r.publicHub, cr)
		r.publicHubMU.Unlock()
		r.logger.Debugln("Event: 'disconnected';", len(r.publicHub), "users")
	}()

	cr <- rr
	r.publicHubMU.RLock()
	r.logger.Debugln("Event: 'connected';", len(r.publicHub), "users")
	r.publicHubMU.RUnlock()

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
