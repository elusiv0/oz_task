package router

import (
	"log/slog"
	"net/http"

	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/elusiv0/oz_task/internal/middleware"
	"github.com/elusiv0/oz_task/internal/router/gql"
	"github.com/elusiv0/oz_task/internal/service"
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	logger *slog.Logger,
	gqlConf graph.Config,
	commentService service.CommentService,
) *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	gql.InitRoutes(logger, router, gqlConf, commentService)

	return router
}
