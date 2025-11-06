package bot

import (
	logging "bookmark-bot/internal"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var commands = []discord.ApplicationCommandCreate{
	discord.MessageCommandCreate{
		Name: "Bookmark local",
	},
}

func registerCommands(client bot.Client) {
	slog.Info("registering commands")
	if err := handler.SyncCommands(client, commands, []snowflake.ID{GuildID}); err != nil {
		logging.FatalLog("error registering commands", err)
	}
}
