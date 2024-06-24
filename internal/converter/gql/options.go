package converter

import (
	"github.com/elusiv0/oz_task/internal/dto"
)

type commentsReqOptions func(*dto.GetCommentsRequest)

func WithCommentsPagination(first int, after *int) commentsReqOptions {
	return func(c *dto.GetCommentsRequest) {
		c.After = after
		c.First = first
	}
}

func WithPostId(postId *int) commentsReqOptions {
	return func(c *dto.GetCommentsRequest) {
		c.PostId = postId
	}
}

func WithParentId(parentId *int) commentsReqOptions {
	return func(c *dto.GetCommentsRequest) {
		c.ParentId = parentId
	}
}

type postsReqOptions func(*dto.GetPostsRequest)

func WithPostsPagination(first int, after *int) postsReqOptions {
	return func(p *dto.GetPostsRequest) {
		p.After = after
		p.First = first
	}
}
