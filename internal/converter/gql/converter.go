package converter

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ToPostConnection(postDto []*dto.Post, first int, after *int) *graph.PostConnection {
	if len(postDto) == 0 {
		return nil
	}
	var edges []*graph.PostEdge
	hasNext := false
	if len(postDto) > first {
		hasNext = true
		postDto = postDto[:len(postDto)-1]
	}
	pageInfo := getPageInfo(postDto[0].ID, postDto[len(postDto)-1].ID, &hasNext)
	for _, val := range postDto {
		curPostEdge := toPostEdge(val)
		edges = append(edges, curPostEdge)
	}

	return &graph.PostConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
}

func toPostEdge(postDto *dto.Post) *graph.PostEdge {
	return &graph.PostEdge{
		Node:   postDto,
		Cursor: postDto.ID,
	}
}

func ToCommentConnection(commentDto []*dto.Comment, first int, after *int) *graph.CommentConnection {
	if len(commentDto) == 0 {
		return nil
	}
	var edges []*graph.CommentEdge
	hasNext := false
	if len(commentDto) > first {
		hasNext = true
		commentDto = commentDto[:len(commentDto)-1]
	}
	pageInfo := getPageInfo(commentDto[0].ID, commentDto[len(commentDto)-1].ID, &hasNext)
	for _, val := range commentDto {
		curCommentEdge := toCommentEdge(val)
		edges = append(edges, curCommentEdge)
	}

	return &graph.CommentConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
}

func toCommentEdge(commentDto *dto.Comment) *graph.CommentEdge {
	return &graph.CommentEdge{
		Node:   commentDto,
		Cursor: commentDto.ID,
	}
}

func getPageInfo(first int, end int, hasNext *bool) *graph.PageInfo {
	return &graph.PageInfo{
		StartCursor: first,
		EndCursor:   end,
		HasNextPage: hasNext,
	}
}

func ToGetCommentsRequest(opts ...commentsReqOptions) dto.GetCommentsRequest {
	commentsReq := &dto.GetCommentsRequest{}
	for _, opt := range opts {
		opt(commentsReq)
	}

	return *commentsReq
}

func ToGetPostsRequest(opts ...postsReqOptions) dto.GetPostsRequest {
	postsReq := &dto.GetPostsRequest{}
	for _, opt := range opts {
		opt(postsReq)
	}

	return *postsReq
}

func ToGqlError(ctx context.Context, cErr *dto.CustomError) *gqlerror.Error {
	return &gqlerror.Error{
		Message: cErr.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"status_code": cErr.GetStatus(),
			"request":     cErr.GetRequestInfo(),
		},
	}
}
