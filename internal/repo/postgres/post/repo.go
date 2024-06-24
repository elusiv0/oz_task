package post

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/repo"
	"github.com/elusiv0/oz_task/internal/repo/converter"
	"github.com/elusiv0/oz_task/internal/repo/model"
	"github.com/elusiv0/oz_task/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

type PostRepository struct {
	db     *postgres.Postgres
	logger *slog.Logger
}

func New(
	postgres *postgres.Postgres,
	logger *slog.Logger,
) *PostRepository {
	repo := &PostRepository{
		db:     postgres,
		logger: logger,
	}

	return repo
}

var _ repo.PostRepo = &PostRepository{}

const (
	postTable = "posts"
)

// Get implements repo.PostRepo.
func (p *PostRepository) Get(ctx context.Context, id int) (*dto.Post, error) {
	postModel := &model.Post{}
	postResp := &dto.Post{}
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - Get - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	sql, args, err := p.db.Builder.
		Select("id", "title", "_text", "closed", "created_at").
		From(postTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - Get - build sql: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&postModel.Id, &postModel.Title,
		&postModel.Text, &postModel.Closed,
		&postModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = dto.NewCustomError(repo.CommentsNotFoundErr, id)
		}
		return postResp, err
	}
	postResp = converter.PostFromRepo(postModel)

	return postResp, nil
}

// GetMany implements repo.PostRepo.
func (p *PostRepository) GetMany(ctx context.Context, postsReq dto.GetPostsRequest) ([]*dto.Post, error) {
	postResp := []*dto.Post{}
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - GetMany - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	builder := p.db.Builder.
		Select("id", "title", "_text", "closed", "created_at").
		From(postTable)
	if postsReq.After != nil {
		builder = builder.Where(squirrel.Lt{"id": postsReq.After})
	}
	sql, args, err := builder.OrderBy("id DESC").
		Limit(uint64(postsReq.First + 1)).
		ToSql()
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - GetMany - build sql: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return postResp, err
	}
	for rows.Next() {
		currPost := &model.Post{}

		err := rows.Scan(
			&currPost.Id, &currPost.Title,
			&currPost.Text, &currPost.Closed,
			&currPost.CreatedAt,
		)
		if err != nil {
			return postResp, fmt.Errorf("PostRepository - GetMany - row scan: %w", err)
		}

		currPostDto := converter.PostFromRepo(currPost)
		postResp = append(postResp, currPostDto)
	}

	if len(postResp) < 2 {
		return postResp, dto.NewCustomError(repo.PostsNotFoundErr, postsReq)
	}

	return postResp, nil
}

// Insert implements repo.PostRepo.
func (p *PostRepository) Insert(ctx context.Context, newPost dto.NewPost) (*dto.Post, error) {
	postResp := &model.Post{}
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return &dto.Post{}, fmt.Errorf("PostRepository - Insert - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	sql, args, err := p.db.Builder.
		Insert(postTable).
		Columns("title", "_text", "closed").
		Values(
			newPost.Title, newPost.Text, newPost.Closed,
		).
		Suffix("RETURNING id, title, _text, closed, created_at").
		ToSql()
	if err != nil {
		return &dto.Post{}, fmt.Errorf("PostRepository - Insert - build sql: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&postResp.Id, &postResp.Title,
		&postResp.Text, &postResp.Closed,
		&postResp.CreatedAt,
	)
	if err != nil {
		return &dto.Post{}, fmt.Errorf("CommentRepository - Insert - scanL %w", err)
	}
	postRespDto := converter.PostFromRepo(postResp)

	return postRespDto, nil
}
