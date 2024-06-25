package post

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/middleware"
	"github.com/elusiv0/oz_task/internal/repo"
	"github.com/elusiv0/oz_task/internal/service"
)

type PostService struct {
	postRepo repo.PostRepo
	logger   *slog.Logger
}

func New(
	postRepo repo.PostRepo,
	logger *slog.Logger,
) *PostService {
	return &PostService{
		postRepo: postRepo,
		logger:   logger,
	}
}

var _ service.PostService = &PostService{}

// Get implements service.PostService.
func (p *PostService) Get(ctx context.Context, id int) (*dto.Post, error) {
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling post repo...")
	postResp, err := p.postRepo.Get(ctx, id)
	if err != nil {
		return postResp, fmt.Errorf("PostService - Get: %w", err)
	}
	logger.Debug("response was handled successfully")

	return postResp, nil
}

// GetMany implements service.PostService.
func (p *PostService) GetMany(ctx context.Context, postsReq dto.GetPostsRequest) ([]*dto.Post, error) {
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling post repo...")
	postResp, err := p.postRepo.GetMany(ctx, postsReq)
	if err != nil {
		return postResp, fmt.Errorf("PostService - GetMany: %w", err)
	}
	logger.Debug("response was handled successfully")

	return postResp, nil
}

// Insert implements service.PostService.
func (p *PostService) Insert(ctx context.Context, newPost dto.NewPost) (*dto.Post, error) {
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("calling post repo...")
	postResp, err := p.postRepo.Insert(ctx, newPost)
	if err != nil {
		return postResp, fmt.Errorf("PostService - Insert: %w", err)
	}
	logger.Debug("response was handled successfully")

	return postResp, nil
}
