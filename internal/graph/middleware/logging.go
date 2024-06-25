package middleware

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/elusiv0/oz_task/internal/middleware"
)

const (
	operationContextKey = "operationContext"
)

func ResponseMiddleware(logger *slog.Logger) graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		resp := next(ctx)
		reqUuid := middleware.GetUuid(ctx)
		logger := logger.With(
			slog.String("request_id", reqUuid),
			slog.Any("path", graphql.GetPath(ctx).String()),
		)
		status := "success"
		level := slog.LevelInfo
		if len(resp.Errors) > 0 {
			status = "error"
			level = slog.LevelWarn
		}
		logger.LogAttrs(ctx, level, "incoming request", slog.String("status", status))
		resp.Extensions = map[string]any{
			"request_id": reqUuid,
		}
		return resp
	}
}
