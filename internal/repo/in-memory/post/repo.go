package post

import (
	"context"
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

type PostRepository struct {
	logger *slog.Logger
	data   map[int]*model.Post
	mu     sync.RWMutex
}

func New(
	logger *slog.Logger,
) *PostRepository {
	return &PostRepository{
		logger: logger,
		data:   make(map[int]*model.Post),
	}
}

var _ repo.PostRepo = &PostRepository{}

var idgen *util.Prid = util.NewPrid()

// Get implements repo.PostRepo.
func (p *PostRepository) Get(ctx context.Context, id int) (*dto.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	postModel, ok := p.data[id]
	if !ok {
		return &dto.Post{}, dto.NewCustomError(repo.PostsNotFoundErr, id)
	}
	postResp := converter.PostFromRepo(postModel)

	return postResp, nil
}

// GetMany implements repo.PostRepo.
func (p *PostRepository) GetMany(ctx context.Context, postsReq dto.GetPostsRequest) ([]*dto.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var posts []*model.Post

	for _, post := range p.data {
		if postsReq.After == nil || post.Id < *postsReq.After {
			posts = append(posts, post)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Id > posts[j].Id
	})
	if len(posts) > postsReq.First {
		posts = posts[:postsReq.First+1]
	}
	var postsDto []*dto.Post
	for _, post := range posts {
		curPostDto := converter.PostFromRepo(post)
		postsDto = append(postsDto, curPostDto)
	}

	if len(postsDto) < 2 {
		return postsDto, dto.NewCustomError(repo.PostsNotFoundErr, postsReq)
	}

	return postsDto, nil
}

// Insert implements repo.PostRepo.
func (p *PostRepository) Insert(ctx context.Context, newPost dto.NewPost) (*dto.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	postModel := &model.Post{
		Id:        idgen.GenerateId(),
		Title:     newPost.Title,
		Text:      newPost.Text,
		Closed:    newPost.Closed,
		CreatedAt: time.Now(),
	}
	p.data[postModel.Id] = postModel
	postResp := converter.PostFromRepo(postModel)

	return postResp, nil
}
