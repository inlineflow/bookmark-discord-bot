package main

import (
	"bookmark-bot/internal/bot"
	"bookmark-bot/internal/logging"
	"log/slog"

	"github.com/disgoorg/disgo"
)

func init() {
	logging.SetDefaultLogger("info")
	slog.Info("starting Bookmark bot")
	slog.Info("disgo version", "version", disgo.Version)
}

func main() {
	bot.Run()
}
