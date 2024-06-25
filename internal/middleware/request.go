package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	reqUuidKey = "reqUuid"
)

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		withCtx := context.WithValue(c.Request.Context(), reqUuidKey, id.String())
		c.Request = c.Request.WithContext(withCtx)
		c.Next()
	}
}

func GetUuid(ctx context.Context) string {
	return ctx.Value(reqUuidKey).(string)
}
