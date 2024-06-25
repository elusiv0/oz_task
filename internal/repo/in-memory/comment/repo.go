package comment

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/repo"
	"github.com/elusiv0/oz_task/internal/repo/converter"
	"github.com/elusiv0/oz_task/internal/repo/model"
	"github.com/elusiv0/oz_task/internal/util"
)

type CommentRepository struct {
	logger *slog.Logger
	data   map[int]*model.Comment
	mu     sync.RWMutex
}

func New(
	logger *slog.Logger,
) *CommentRepository {
	return &CommentRepository{
		logger: logger,
		data:   make(map[int]*model.Comment),
	}
}

var _ repo.CommentRepo = &CommentRepository{}

var idgen *util.Prid = util.NewPrid()

// Get implements repo.CommentRepo.
func (c *CommentRepository) Get(ctx context.Context, id int) (*dto.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	commentModel, ok := c.data[id]
	if !ok {
		return &dto.Comment{}, dto.NewCustomError(repo.CommentsNotFoundErr, id)
	}
	commentResp := converter.CommentFromRepo(commentModel)

	return commentResp, nil
}

// GetMany implements repo.CommentRepo.
func (c *CommentRepository) GetMany(ctx context.Context, commentsReq ...dto.GetCommentsRequest) ([]*dto.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	commentssl := make(map[int][]*model.Comment)
	set := make(map[int]bool)
	byParent := false
	byPost := false
	after := commentsReq[0].After
	first := commentsReq[0].First
	if commentsReq[0].ParentId != nil {
		byParent = true
	} else if commentsReq[0].PostId != nil {
		byPost = true
	}
	if !byPost && !byParent {
		return []*dto.Comment{}, fmt.Errorf("CommentRepository - GetMany: required postId or parentId")
	}
	for _, val := range commentsReq {
		if val.ParentId != nil {
			set[*val.ParentId] = true
		} else {
			set[*val.PostId] = true
		}
	}

	for _, post := range c.data {
		if byParent {
			if !post.ParentId.Valid {
				continue
			}
			parentId := int(post.ParentId.Int32)
			if _, ok := set[parentId]; ok {
				if after == nil || post.Id < *after {
					commentssl[parentId] = append(commentssl[parentId], post)
				}
			}
		} else {
			postId := post.ArticleID
			if _, ok := set[postId]; ok && !post.ParentId.Valid {
				if after == nil || post.Id < *after {
					commentssl[postId] = append(commentssl[postId], post)
				}
			}
		}
	}
	for _, val := range commentssl {
		sort.Slice(val, func(i, j int) bool {
			return val[i].Id > val[j].Id
		})
	}
	var commentsResp []*dto.Comment

	for _, comm := range commentssl {
		if len(comm) > first {
			comm = comm[:first+1]
		}
		for _, commBy := range comm {
			currCommentDto := converter.CommentFromRepo(commBy)
			commentsResp = append(commentsResp, currCommentDto)
		}
	}

	if len(commentsResp) == 0 && len(commentsReq) == 1 {
		return []*dto.Comment{}, dto.NewCustomError(repo.CommentsNotFoundErr, commentsReq[0])
	}
	return commentsResp, nil
}

// Insert implements repo.CommentRepo.
func (c *CommentRepository) Insert(ctx context.Context, newComment dto.NewComment) (*dto.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	pId := sql.NullInt32{
		Valid: false,
	}
	if newComment.ParentID != nil {
		pId.Int32 = int32(*newComment.ParentID)
		pId.Valid = true
	}
	commentModel := &model.Comment{
		Id:        idgen.GenerateId(),
		Text:      newComment.Text,
		ArticleID: newComment.ArticleID,
		ParentId:  pId,
		CreatedAt: time.Now(),
	}
	c.data[commentModel.Id] = commentModel
	commentResp := converter.CommentFromRepo(commentModel)
	return commentResp, nil
}
