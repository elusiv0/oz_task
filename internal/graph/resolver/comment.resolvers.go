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
	"github.com/elusiv0/oz_task/internal/graph/dataloader"
)

// Comments is the resolver for the comments field.
func (r *commentResolver) Comments(ctx context.Context, obj *model.Comment, first *int, after *int) (*graph.CommentConnection, error) {
	commentsReq := gqlconv.ToGetCommentsRequest(
		gqlconv.WithCommentsPagination(*first, after),
		gqlconv.WithParentId(&obj.ID),
	)
	commentsResp, err := dataloader.GetCommentLoader(ctx).Load(commentsReq)
	if err != nil {
		r.logger.Warn("Error was handled", slog.String("Cause", "PostResolver - Comments: "+err.Error()))
		gqlErr := handleError(ctx, err)
		return &graph.CommentConnection{}, gqlErr
	}
	commentConn := gqlconv.ToCommentConnection(commentsResp, commentsReq.First, commentsReq.After)

	return commentConn, nil
}

// Comment returns graph.CommentResolver implementation.
func (r *Resolver) Comment() graph.CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }
