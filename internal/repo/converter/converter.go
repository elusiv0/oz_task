package converter

import (
	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/repo/model"
)

func CommentToRepo(commentDto *dto.Comment) *model.Comment {
	return &model.Comment{
		Id:        commentDto.ID,
		Text:      commentDto.Text,
		ArticleID: commentDto.ArticleID,
		CreatedAt: commentDto.CreatedAt,
	}
}

func CommentFromRepo(commentModel *model.Comment) *dto.Comment {
	var pId *int
	if commentModel.ParentId.Valid {
		elem := int(commentModel.ParentId.Int32)
		pId = &elem
	}
	return &dto.Comment{
		ID:        commentModel.Id,
		Text:      commentModel.Text,
		ArticleID: commentModel.ArticleID,
		ParentID:  pId,
		CreatedAt: commentModel.CreatedAt,
	}
}

func PostFromRepo(postModel *model.Post) *dto.Post {
	return &dto.Post{
		ID:        postModel.Id,
		Text:      postModel.Text,
		Title:     postModel.Title,
		Closed:    postModel.Closed,
		CreatedAt: postModel.CreatedAt,
	}
}

func PostToRepo(postDto *dto.Post) *model.Post {
	return &model.Post{
		Id:        postDto.ID,
		Text:      postDto.Text,
		Title:     postDto.Title,
		Closed:    postDto.Closed,
		CreatedAt: postDto.CreatedAt,
	}
}
