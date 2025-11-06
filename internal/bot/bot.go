package bot

import (
	logging "bookmark-bot/internal"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

func setup() {
	if err := godotenv.Load(); err != nil {
		logging.FatalLog("error loading .env", err)
	}

	Token = os.Getenv("DISCORD_BOT_TOKEN")
	rawGuildID := os.Getenv("DISCORD_GUILD_ID")
	GuildID = snowflake.MustParse(rawGuildID)
}

func Run() {
	setup()

	b := newBot()

	client, err := disgo.New(Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages, gateway.IntentDirectMessages),
		),
	)
	if err != nil {
		slog.Error("error starting disgo client", "error", err)
		os.Exit(1)
	}

	b.Client = client

	registerCommands(client)

	b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.MessageCommandInteractionData) error{
		"Bookmark": b.bookmark,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		slog.Error("failed to open gateway", "error", err)
		os.Exit(1)
	}
	defer client.Close(context.TODO())

	slog.Info("bot is now running, press CTRL-C to exit")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}

func newBot() *Bot {
	return &Bot{}
}

type Bot struct {
	Client   bot.Client
	Handlers map[string]func(e *events.ApplicationCommandInteractionCreate, d discord.MessageCommandInteractionData) error
}
