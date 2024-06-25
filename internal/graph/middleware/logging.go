package middleware

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

const (
	reqUuidKey          = "reqUuid"
	operationContextKey = "operationContext"
)

func ReqUuuidMiddleware() graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		reqUuid := uuid.New()
		oc := graphql.GetOperationContext(ctx)
		ctxWith := context.WithValue(ctx, reqUuidKey, reqUuid.String())
		ctxWith = context.WithValue(ctxWith, operationContextKey, oc)
		return next(ctxWith)
	}
}

func GetReqUuid(ctx context.Context) string {
	reqUuid := ctx.Value(reqUuidKey).(string)
	return reqUuid
}

func getOperationContext(ctx context.Context) *graphql.OperationContext {
	oc := ctx.Value(operationContextKey).(*graphql.OperationContext)
	return oc
}

func LoggingMiddleware(logger *slog.Logger) graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		resp := next(ctx)
		oc := getOperationContext(ctx)
		reqUuid := GetReqUuid(ctx)
		logger := logger.With(
			slog.String("request", reqUuid),
			slog.String("operation_name", oc.OperationName),
			slog.Any("path", graphql.GetPath(ctx).String()),
		)
		status := "success"
		level := slog.LevelInfo
		if len(resp.Errors) > 0 {
			status = "error"
			level = slog.LevelWarn
		}
		logger.LogAttrs(ctx, level, "", slog.String("status", status))
		return resp
	}
}
