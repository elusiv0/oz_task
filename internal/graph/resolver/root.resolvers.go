package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"log/slog"

	gqlconv "github.com/elusiv0/oz_task/internal/converter/gql"
	model "github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/google/uuid"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	postResp, err := r.postService.Insert(ctx, input)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - CreatePost: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return postResp, gqlErr
	}

	return postResp, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	commentResp, err := r.commentService.Insert(ctx, input)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - CreateComment: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return commentResp, gqlErr
	}

	for _, ch := range r.postsSubscribers[commentResp.ArticleID] {
		ch <- commentResp
	}

	return commentResp, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, first *int, after *int) (*graph.PostConnection, error) {
	postsReq := gqlconv.ToGetPostsRequest(
		gqlconv.WithPostsPagination(*first, after),
	)
	postResp, err := r.postService.GetMany(ctx, postsReq)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Posts: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return &graph.PostConnection{}, gqlErr
	}

	postConn := gqlconv.ToPostConnection(postResp, postsReq.First, postsReq.After)

	return postConn, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id *int) (*model.Post, error) {
	postResp, err := r.postService.Get(ctx, *id)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Post: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return postResp, gqlErr
	}

	return postResp, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id *int) (*model.Comment, error) {
	commentResp, err := r.commentService.Get(ctx, *id)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "mutationResolver - Post: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return commentResp, gqlErr
	}

	return commentResp, nil
}

// NewComments is the resolver for the newComments field.
func (r *subscriptionResolver) NewComments(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	id := uuid.New()
	comments := make(chan *model.Comment, 1)

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
