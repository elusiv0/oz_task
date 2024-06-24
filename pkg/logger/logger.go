package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func New(env string) (log *slog.Logger) {

	switch env {
	case local:
		log = slog.New(
			tint.NewHandler(
				colorable.NewColorable(os.Stdout),
				&tint.Options{
					Level: slog.LevelDebug,
				},
			),
		)
	case dev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)
	case prod:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	}

	return
}
