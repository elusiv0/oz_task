package resolver

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	gqlconv "github.com/elusiv0/oz_task/internal/converter/gql"
	model "github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/service"
	"github.com/elusiv0/oz_task/internal/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Resolver struct {
	commentService   service.CommentService
	postService      service.PostService
	logger           *slog.Logger
	postsSubscribers map[int]map[string]chan *model.Comment
	mu               sync.Mutex
}

var customError *model.CustomError

func NewResolver(
	commentService service.CommentService,
	postService service.PostService,
	logger *slog.Logger,
) *Resolver {
	return &Resolver{
		logger:           logger,
		commentService:   commentService,
		postService:      postService,
		postsSubscribers: make(map[int]map[string]chan *model.Comment),
	}
}

func handleError(ctx context.Context, err error) error {
	cause := util.UnwrapError(err)
	gqlErr := &gqlerror.Error{
		Message: "Internal Server Error",
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]any{
			"status": 500,
		},
	}
	if errors.As(cause, &customError) {
		cErr := cause.(*model.CustomError)
		gqlErr = gqlconv.ToGqlError(ctx, cErr)
	}

	return gqlErr
}
