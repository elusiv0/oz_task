package route

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func InitRoutes(
	logger *slog.Logger,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(sloggin.New(logger))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	//_ := router.Group("/api/v1/orders/")
	{
		//orderRouter.NewRouter(orderGroup, cache, logger)
	}

	return router
}
