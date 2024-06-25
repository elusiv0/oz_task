package post

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/middleware"
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
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("initialize transaction...")
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - Get - begin tx: %w", err)
	}
	logger.Debug("transation was initialized successfully")
	defer func() {
		if err != nil {
			logger.Debug("error was handled, rollback transaction")
			tx.Rollback(ctx)
			return
		}
		logger.Debug("transaction was committed successfully")
		err = tx.Commit(ctx)
	}()

	logger.Debug("building sql...")
	sql, args, err := p.db.Builder.
		Select("id", "title", "_text", "closed", "created_at").
		From(postTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - Get - build sql: %w", err)
	}
	logger.Debug("sql was builded successfully", slog.String("sql", sql), slog.Any("args", args))

	logger.Debug("executing sql statement...")
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
	logger.Debug("sql statement was executed successfully")

	logger.Debug("converting comment model to dto...")
	postResp = converter.PostFromRepo(postModel)
	logger.Debug("model was converted successfully")

	return postResp, nil
}

// GetMany implements repo.PostRepo.
func (p *PostRepository) GetMany(ctx context.Context, postsReq dto.GetPostsRequest) ([]*dto.Post, error) {
	postResp := []*dto.Post{}
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("initialize transaction...")
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return postResp, fmt.Errorf("PostRepository - GetMany - begin tx: %w", err)
	}
	logger.Debug("transation was initialized successfully")
	defer func() {
		if err != nil {
			logger.Debug("error was handled, rollback transaction")
			tx.Rollback(ctx)
			return
		}
		logger.Debug("transaction was committed successfully")
		err = tx.Commit(ctx)
	}()

	logger.Debug("building sql...")
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
	logger.Debug("sql was builded successfully", slog.String("sql", sql), slog.Any("args", args))

	logger.Debug("executing sql statement...")
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return postResp, err
	}
	defer rows.Close()
	logger.Debug("sql statement was executed successfully")

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

		logger.Debug("converting comment model to dto...")
		currPostDto := converter.PostFromRepo(currPost)
		logger.Debug("model was converted successfully")
		postResp = append(postResp, currPostDto)
	}

	if len(postResp) == 0 {
		return postResp, dto.NewCustomError(repo.PostsNotFoundErr, postsReq)
	}

	return postResp, nil
}

// Insert implements repo.PostRepo.
func (p *PostRepository) Insert(ctx context.Context, newPost dto.NewPost) (*dto.Post, error) {
	postResp := &model.Post{}
	logger := p.logger.With(slog.String("request_id", middleware.GetUuid(ctx)))

	logger.Debug("initialize transaction...")
	tx, err := p.db.PgxPool.Begin(ctx)
	if err != nil {
		return &dto.Post{}, fmt.Errorf("PostRepository - Insert - begin tx: %w", err)
	}
	logger.Debug("transation was initialized successfully")
	defer func() {
		if err != nil {
			logger.Debug("error was handled, rollback transaction")
			tx.Rollback(ctx)
			return
		}
		logger.Debug("transaction was committed successfully")
		err = tx.Commit(ctx)
	}()

	logger.Debug("building sql...")
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
	logger.Debug("sql was builded successfully", slog.String("sql", sql), slog.Any("args", args))

	logger.Debug("executing sql statement...")
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&postResp.Id, &postResp.Title,
		&postResp.Text, &postResp.Closed,
		&postResp.CreatedAt,
	)
	if err != nil {
		return &dto.Post{}, fmt.Errorf("CommentRepository - Insert - scanL %w", err)
	}
	logger.Debug("sql statement was executed successfully")

	logger.Debug("converting comment model to dto...")
	postRespDto := converter.PostFromRepo(postResp)
	logger.Debug("model was converted successfully")

	return postRespDto, nil
}
