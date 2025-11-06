package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (b *Bot) bookmark(e *events.ApplicationCommandInteractionCreate, d discord.MessageCommandInteractionData) error {
	msg := d.Resolved.Messages[d.TargetID()]
	ptr := &msg
	if ptr == nil {
		return InvisibleReply("Could not find the message.", e)
	}

	dmChannel, err := b.Client.Rest().CreateDMChannel(e.User().ID)
	if err != nil {
		slog.Error("failed to create DM channel", "error", err)
	}
	msgChannel, err := b.Client.Rest().GetChannel(e.Channel().ID())
	if err != nil {
		slog.Error("failed to fetch channel", "error", err)
	}
	msgGuild, err := b.Client.Rest().GetChannel(*e.GuildID())
	if err != nil {
		slog.Error("failed to fetch guild", "error", err)
	}

	return EmbedResponse(b, msg, dmChannel, msgChannel, msgGuild, e)
}
