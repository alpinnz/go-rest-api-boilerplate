package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	// Pilih handler berdasarkan ENV
	// Kalau DEV → pakai TextHandler biar lebih kebaca
	// Kalau PROD → pakai JSONHandler biar bisa di-parse monitoring tools
	env := os.Getenv("APP_ENV")

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // default: INFO ke atas
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug, // di DEV log semua level
		})
	}

	Log = slog.New(handler)
}
