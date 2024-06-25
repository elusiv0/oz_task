package resolver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	gqlconv "github.com/elusiv0/oz_task/internal/converter/gql"
	model "github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/service"
	"github.com/elusiv0/oz_task/internal/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	CreateCommentPostNotFound = model.ErrInfo{
		ErrorMessage: "couldn't get required post with provided id",
		StatusCode:   http.StatusBadRequest,
	}
	CreateCommentPostClosedErr = model.ErrInfo{
		ErrorMessage: "post closed to add comments",
		StatusCode:   http.StatusForbidden,
	}
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
		if cErr.GetStatus() == http.StatusNoContent {
			return nil
		}
		gqlErr = gqlconv.ToGqlError(ctx, cErr)
	}

	return gqlErr
}
