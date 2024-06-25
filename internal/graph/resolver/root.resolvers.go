package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"log/slog"

	gqlconv "github.com/elusiv0/oz_task/internal/converter/gql"
	model "github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/elusiv0/oz_task/internal/middleware"
	"github.com/google/uuid"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling comment service...")
	postResp, err := r.postService.Insert(ctx, input)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - CreatePost: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	return postResp, nil
}

// CreateComment implements graph.MutationResolver.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling post service...")
	post, err := r.postService.Get(ctx, input.ArticleID)
	if err != nil {
		logger.Warn("Error was handled", slog.String("Cause", "PostResolver - Comments: "+err.Error()))
		var customErr *model.CustomError
		if errors.As(err, &customErr) {
			err = model.NewCustomError(CreateCommentPostNotFound, input)
		}
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}
	if post.Closed {
		err := model.NewCustomError(CreateCommentPostClosedErr, input)
		logger.Warn("Error was handled", slog.String("Cause", "PostResolver - Comments: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	logger.Debug("calling comment service...")
	commentResp, err := r.commentService.Insert(ctx, input)
	if err != nil {
		logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - CreateComment: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	logger.Debug("sending comment response to subscribe channel...")
	for _, ch := range r.postsSubscribers[commentResp.ArticleID] {
		ch <- commentResp
	}

	return commentResp, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, first *int, after *int) (*graph.PostConnection, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("wrapping post request to dto...")
	postsReq := gqlconv.ToGetPostsRequest(
		gqlconv.WithPostsPagination(*first, after),
	)

	logger.Debug("calling post service...")
	postResp, err := r.postService.GetMany(ctx, postsReq)
	if err != nil {
		logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Posts: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	logger.Debug("converting post response to post connection...")
	postConn := gqlconv.ToPostConnection(postResp, postsReq.First, postsReq.After)

	return postConn, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id *int) (*model.Post, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling post service...")
	postResp, err := r.postService.Get(ctx, *id)
	if err != nil {
		logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Post: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	return postResp, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id *int) (*model.Comment, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling comment service...")
	commentResp, err := r.commentService.Get(ctx, *id)
	if err != nil {
		logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Post: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return nil, gqlErr
	}

	return commentResp, nil
}

// NewComments is the resolver for the newComments field.
func (r *subscriptionResolver) NewComments(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	logger := r.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))
	id := uuid.New()
	comments := make(chan *model.Comment, 1)

	logger.Debug("initializing stop subscription goroutine...")
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.postsSubscribers[postID], id.String())
		r.mu.Unlock()
	}()

	if r.postsSubscribers[postID] == nil {
		r.postsSubscribers[postID] = make(map[string]chan *model.Comment)
	}

	r.postsSubscribers[postID][id.String()] = comments
	logger.Debug("subscription is ready")

	return comments, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
