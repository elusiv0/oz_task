package main

import (
	"log"

	"github.com/elusiv0/oz_task/internal/app"
	"github.com/elusiv0/oz_task/internal/config"
	"github.com/elusiv0/oz_task/internal/graph"
	resolver "github.com/elusiv0/oz_task/internal/graph/resolver"
	"github.com/elusiv0/oz_task/internal/repo"
	imCommentRepo "github.com/elusiv0/oz_task/internal/repo/in-memory/comment"
	imPostRepo "github.com/elusiv0/oz_task/internal/repo/in-memory/post"
	pgCommentRepo "github.com/elusiv0/oz_task/internal/repo/postgres/comment"
	pgPostRepo "github.com/elusiv0/oz_task/internal/repo/postgres/post"
	"github.com/elusiv0/oz_task/internal/router"
	commentService "github.com/elusiv0/oz_task/internal/service/comment"
	postService "github.com/elusiv0/oz_task/internal/service/post"
	"github.com/elusiv0/oz_task/pkg/httpserver"
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

	//building repo
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

	//building service
	commentService := commentService.New(commentRepo, logger)
	postService := postService.New(postRepo, logger)

	//building gql
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

	//building router
	router := router.InitRoutes(logger, gConfig, commentService)

	//building httpserver
	httpserver := httpserver.New(
		router,
		httpserver.Port(config.Http.Port),
		httpserver.ReadTimeout(config.Http.ReadTimeout),
		httpserver.ShutdownTimeout(config.Http.ShutdownTimeout),
	)

	//building app
	app := app.New(httpserver, logger)

	logger.Info("Starting app on on port" + config.Http.Port + "...")
	if err := app.Run(); err != nil {
		log.Fatal("error with starting app" + err.Error())
	}
}
