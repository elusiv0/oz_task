package service

import (
	"context"

	"github.com/elusiv0/oz_task/internal/dto"
)

type PostService interface {
	GetMany(ctx context.Context, postsReq dto.GetPostsRequest) ([]*dto.Post, error)
	Insert(ctx context.Context, newPost dto.NewPost) (*dto.Post, error)
	Get(ctx context.Context, id int) (*dto.Post, error)
}

type CommentService interface {
	GetMany(ctx context.Context, commentsReq ...dto.GetCommentsRequest) ([]*dto.Comment, error)
	Insert(ctx context.Context, newComment dto.NewComment) (*dto.Comment, error)
	Get(ctx context.Context, id int) (*dto.Comment, error)
}
