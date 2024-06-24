package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/elusiv0/oz_task/internal/config"
	"github.com/elusiv0/oz_task/internal/graph"
	"github.com/elusiv0/oz_task/internal/graph/dataloader"
	resolver "github.com/elusiv0/oz_task/internal/graph/resolver"
	"github.com/elusiv0/oz_task/internal/repo"
	imCommentRepo "github.com/elusiv0/oz_task/internal/repo/in-memory/comment"
	imPostRepo "github.com/elusiv0/oz_task/internal/repo/in-memory/post"
	pgCommentRepo "github.com/elusiv0/oz_task/internal/repo/postgres/comment"
	pgPostRepo "github.com/elusiv0/oz_task/internal/repo/postgres/post"
	commentService "github.com/elusiv0/oz_task/internal/service/comment"
	postService "github.com/elusiv0/oz_task/internal/service/post"
	"github.com/elusiv0/oz_task/pkg/logger"
	"github.com/elusiv0/oz_task/pkg/postgres"
	"github.com/joho/godotenv"
)

func main() {
	//load env variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error with load env variables " + err.Error())
	}

	//building config
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal("error with get config " + err.Error())
	}

	//building logger
	logger := logger.New("local")

	var postRepo repo.PostRepo
	var commentRepo repo.CommentRepo
	if config.App.Db == "postgres" {
		pg, err := postgres.New(
			postgres.NewConnectionConfig(
				config.Postgres.Host,
				config.Postgres.Port,
				config.Postgres.User,
				config.Postgres.Password,
				config.Postgres.Name,
			),
			logger,
			postgres.ConnAttempts(config.Postgres.ConnectionAttempts),
			postgres.ConnTimeout(config.Postgres.ConnectionTimeout),
			postgres.MaxPoolSz(config.Postgres.MaxPoolSz),
		)
		if err != nil {
			log.Fatal("error with set up pg connection " + err.Error())
		}

		postRepo = pgPostRepo.New(pg, logger)
		commentRepo = pgCommentRepo.New(pg, logger)
	} else {
		postRepo = imPostRepo.New(logger)
		commentRepo = imCommentRepo.New(logger)
	}

	commentService := commentService.New(commentRepo, logger)
	postService := postService.New(postRepo, logger)

	resolver := resolver.NewResolver(commentService, postService, logger)
	gConfig := graph.Config{
		Resolvers: resolver,
	}
	countComplexity := func(childComplexity int, first, after *int) int {
		return *first * childComplexity
	}
	gConfig.Complexity.Post.Comments = countComplexity
	gConfig.Complexity.Query.Posts = countComplexity
	gConfig.Complexity.Comment.Comments = countComplexity
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gConfig))

	srv.Use(extension.FixedComplexityLimit(1500))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", dataloader.DataloaderMiddleware(commentService, srv))
	port := "8080"
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	logger.Info("Starting app...")
	//if err := app.Run(); err != nil {
	//	log.Fatal("error with starting app" + err.Error())
	//}

}
