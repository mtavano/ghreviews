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
	publicHub   map[chan *GhReviewsEvent]bool

	privateHubMU sync.RWMutex
	privateHub   map[string][](chan *GhReviewsEvent)
}

func NewResolver(logger *logrus.Logger, reviewService ghreviews.ReviewService) *Resolver {
	return &Resolver{
		reviewService: reviewService,
		logger:        logger,
		publicHub:     make(map[chan *GhReviewsEvent]bool),
		privateHub:    make(map[string][]chan *GhReviewsEvent),
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

	total, err := r.reviewService.CountReviewsByUsername(reviewInput.GithubUsername)
	if err != nil {
		return nil, err
	}

	msg := &GhReviewsEvent{Total: total, NewReviews: []*ghreviews.GhReview{review}}
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

func (r *subscriptionResolver) FeedByUsername(ctx context.Context, username string) (<-chan *GhReviewsEvent, error) {
	rr, err := r.reviewService.GetLastReviewsByUsername(username)
	if err != nil {
		return nil, err
	}
	total, err := r.reviewService.CountReviewsByUsername(username)
	if err != nil {
		return nil, err
	}

	// Add channel to private hub (by username)
	cr := make(chan *GhReviewsEvent, 1)
	r.privateHubMU.Lock()
	r.privateHub[username] = append(r.privateHub[username], cr)
	r.privateHubMU.Unlock()
	go func() {
		<-ctx.Done()
		r.privateHubMU.Lock()
		newList := []chan *GhReviewsEvent{}
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

	cr <- &GhReviewsEvent{Total: total, NewReviews: rr}
	r.publicHubMU.RLock()
	r.logger.Debugf("Event: 'connected to %q'; %d users", username, len(r.privateHub[username]))
	r.publicHubMU.RUnlock()

	return cr, nil
}

func (r *subscriptionResolver) Feed(ctx context.Context) (<-chan *GhReviewsEvent, error) {
	rr, err := r.reviewService.GetLastReviews()
	if err != nil {
		return nil, err
	}

	total, err := r.reviewService.CountReviews()
	if err != nil {
		return nil, err
	}

	// Add channel to hub
	cr := make(chan *GhReviewsEvent, 1)
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

	cr <- &GhReviewsEvent{Total: total, NewReviews: rr}
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
