package gql

import (
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/elusiv0/oz_task/internal/graph/middleware"
	"github.com/elusiv0/oz_task/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	reqUuidKey = "reqUuid"
)

func InitRoutes(
	logger *slog.Logger,
	router *gin.Engine,
	graphConfig graph.Config,
	commentService service.CommentService,
) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graphConfig))
	srv.Use(extension.FixedComplexityLimit(1500))
	srv.AroundResponses(middleware.ResponseMiddleware(logger))
	router.GET("/", playgroundHandler(playground.Handler("GraphQL playground", "/query")))
	router.Any("/query", graphqlHandler(middleware.DataloaderMiddleware(commentService, srv)))
}

func graphqlHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
