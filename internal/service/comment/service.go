package comment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/repo"
	"github.com/elusiv0/oz_task/internal/service"
)

type CommentService struct {
	commentRepo repo.CommentRepo
	logger      *slog.Logger
}

func New(
	commentRepo repo.CommentRepo,
	logger *slog.Logger,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		logger:      logger,
	}
}

var _ service.CommentService = &CommentService{}

// Get implements service.CommentRepo.
func (c *CommentService) Get(ctx context.Context, id int) (*dto.Comment, error) {
	commentResp, err := c.commentRepo.Get(ctx, id)
	if err != nil {
		return commentResp, fmt.Errorf("CommentService - Get: %w", err)
	}

	return commentResp, nil
}

// GetMany implements service.CommentRepo.
func (c *CommentService) GetMany(ctx context.Context, commentReq ...dto.GetCommentsRequest) ([]*dto.Comment, error) {
	commentResp, err := c.commentRepo.GetMany(ctx, commentReq...)
	if err != nil {
		return commentResp, fmt.Errorf("CommentService - GetMany: %w", err)
	}

	return commentResp, nil
}

// Insert implements service.CommentRepo.
func (c *CommentService) Insert(ctx context.Context, newComment dto.NewComment) (*dto.Comment, error) {
	commentResp, err := c.commentRepo.Insert(ctx, newComment)
	if err != nil {
		return commentResp, fmt.Errorf("CommentService - Insert: %w", err)
	}

	return commentResp, nil
}
