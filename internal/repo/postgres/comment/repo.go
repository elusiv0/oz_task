package comment

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

type CommentRepository struct {
	db     *postgres.Postgres
	logger *slog.Logger
}

func New(
	postgres *postgres.Postgres,
	logger *slog.Logger,
) *CommentRepository {
	repo := &CommentRepository{
		db:     postgres,
		logger: logger,
	}

	return repo
}

var _ repo.CommentRepo = &CommentRepository{}

const (
	commentTable = "comments"
)

// Get implements repo.CommentRepo.
func (c *CommentRepository) Get(ctx context.Context, id int) (*dto.Comment, error) {
	commentModel := &model.Comment{}
	commentResp := &dto.Comment{}
	tx, err := c.db.PgxPool.Begin(ctx)
	if err != nil {
		return commentResp, fmt.Errorf("CommentRepository - Get - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	sql, args, err := c.db.Builder.
		Select("id", "_text", "article_id", "parent_id", "created_at").
		From(commentTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return commentResp, fmt.Errorf("CommentRepository - Get - build sql: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&commentModel.Id, &commentModel.Text,
		&commentModel.ArticleID, &commentModel.ParentId,
		&commentModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			cErr := dto.NewCustomError(repo.CommentsNotFoundErr, id)
			err = cErr
		}
		return commentResp, err
	}
	commentResp = converter.CommentFromRepo(commentModel)

	return commentResp, nil
}

// GetMany implements repo.CommentRepo.
func (c *CommentRepository) GetMany(ctx context.Context, commentsReq ...dto.GetCommentsRequest) ([]*dto.Comment, error) {
	commentResp := []*dto.Comment{}
	commentModel := &model.Comment{}
	tx, err := c.db.PgxPool.Begin(ctx)
	if err != nil {
		return commentResp, fmt.Errorf("CommentRepository - GetMany - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()
	var builder *squirrel.SelectBuilder
	var scanRows []any
	if len(commentsReq) > 1 {
		builder, scanRows = c.buildManyVariadic(commentModel, commentsReq...)
	} else {
		builder, scanRows = c.buildManyNonVariadic(commentModel, commentsReq[0])
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return commentResp, fmt.Errorf("CommentRepository - GetMany - build sql: %w", err)
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return commentResp, err
	}

	for rows.Next() {
		err := rows.Scan(
			scanRows...,
		)
		if err != nil {
			return commentResp, fmt.Errorf("CommentRepository - GetMany - row scan: %w", err)
		}
		curCommentDto := converter.CommentFromRepo(commentModel)

		commentResp = append(commentResp, curCommentDto)
	}

	if len(commentResp) < 2 && len(commentsReq) == 1 {
		cErr := dto.NewCustomError(repo.CommentsNotFoundErr, commentsReq[0])

		return commentResp, cErr
	}
	return commentResp, nil
}

// Insert implements repo.CommentRepo.
func (c *CommentRepository) Insert(ctx context.Context, newComment dto.NewComment) (*dto.Comment, error) {
	commentResp := &model.Comment{}
	tx, err := c.db.PgxPool.Begin(ctx)
	if err != nil {
		return &dto.Comment{}, fmt.Errorf("CommentRepository - Insert - begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	sql, args, err := c.db.Builder.
		Insert(commentTable).
		Columns("_text", "article_id", "parent_id").
		Values(
			newComment.Text, newComment.ArticleID, newComment.ParentID,
		).
		Suffix("RETURNING id, _text, article_id, parent_id, created_at").
		ToSql()
	if err != nil {
		return &dto.Comment{}, fmt.Errorf("CommentRepository - Insert - build sql: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&commentResp.Id, &commentResp.Text,
		&commentResp.ArticleID, &commentResp.ParentId,
		&commentResp.CreatedAt,
	)
	if err != nil {
		return &dto.Comment{}, fmt.Errorf("CommentRepository - Insert - scanL %w", err)
	}
	commentRespDto := converter.CommentFromRepo(commentResp)

	return commentRespDto, nil
}

func (c *CommentRepository) buildManyNonVariadic(commentResp *model.Comment, commentsReq dto.GetCommentsRequest) (*squirrel.SelectBuilder, []any) {
	var conditions squirrel.And
	scanRows := []any{
		&commentResp.Id, &commentResp.Text,
		&commentResp.ArticleID, &commentResp.ParentId,
		&commentResp.CreatedAt,
	}
	conditions = append(conditions, squirrel.Eq{"parent_id": commentsReq.ParentId})
	if commentsReq.PostId != nil {
		conditions = append(conditions, squirrel.Eq{"article_id": commentsReq.PostId})
	}
	if commentsReq.After != nil {
		conditions = append(conditions, squirrel.Lt{"id": commentsReq.After})
	}
	builder := c.db.Builder.
		Select("id", "_text", "article_id", "parent_id", "created_at").
		From(commentTable)
	if len(conditions) > 0 {
		builder = builder.Where(conditions)
	}
	builder = builder.
		OrderBy("id DESC").
		Limit(uint64(commentsReq.First + 1))
	return &builder, scanRows
}

func (c *CommentRepository) buildManyVariadic(commentResp *model.Comment, commentsReq ...dto.GetCommentsRequest) (*squirrel.SelectBuilder, []any) {
	var parentsId []*int
	var postsId []*int
	scanRows := []any{
		&commentResp.Id, &commentResp.Text,
		&commentResp.ArticleID, &commentResp.ParentId,
		&commentResp.CreatedAt, &commentResp.Rown,
	}
	partition := "parent_id"
	if commentsReq[0].ParentId != nil {
		parentsId = make([]*int, 10)
	} else {
		postsId = make([]*int, 10)
		partition = "article_id"
	}
	first := commentsReq[0].First

	for _, val := range commentsReq {
		if parentsId != nil {
			parentsId = append(parentsId, val.ParentId)
		}
		if postsId != nil {
			postsId = append(postsId, val.PostId)
		}
	}
	var conditions squirrel.And
	if parentsId != nil {
		conditions = append(conditions, squirrel.Eq{"parent_id": parentsId})
	} else {
		conditions = append(conditions, squirrel.Eq{"parent_id": nil})
	}
	if postsId != nil {
		conditions = append(conditions, squirrel.Eq{"article_id": postsId})
	}
	subSelect := c.db.Builder.
		Select("id", "_text",
			"article_id", "parent_id",
			"created_at", "row_number() OVER (PARTITION BY "+partition+" ORDER BY id DESC) AS com_row").
		From(commentTable).
		Where(conditions)
	builder := c.db.Builder.
		Select("id", "_text", "article_id", "parent_id", "created_at", "com_row").
		FromSelect(subSelect, "com").
		Where(squirrel.LtOrEq{"com.com_row": first})
	return &builder, scanRows
}