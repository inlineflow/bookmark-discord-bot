package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var (
	Token   string
	GuildID snowflake.ID
)

func InvisibleReply(content string, e *events.ApplicationCommandInteractionCreate) error {
	if err := e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(content).
		SetEphemeral(true).
		Build(),
	); err != nil {
		slog.Error("error sending response", "error", err)
		return errors.New("error sending response")
	}
	return nil
}

func EmbedResponse(b *Bot, m discord.Message, dm *discord.DMChannel, c discord.Channel, g discord.Channel, e *events.ApplicationCommandInteractionCreate) error {
	fields := []discord.EmbedField{
		{
			Name:  "From",
			Value: fmt.Sprintf("%s, $%s", g, c),
		},
		{
			Value: fmt.Sprintf("[Link to the original message](https://discord.com/channels/%v/%v/%v)", g.ID(), m.ChannelID, m.ID),
		},
		{
			Name:  "Message preview",
			Value: m.Content,
		},
	}
	if a := len(m.Attachments); a != 0 {
		fields = append(fields, discord.EmbedField{Value: fmt.Sprintf("%d %s", a, pluralize("attachment", a))})
	}
	if e := len(m.Embeds); e != 0 {
		fields = append(fields, discord.EmbedField{Value: fmt.Sprintf("%d %s", e, pluralize("embed", e))})
	}
	embed := discord.NewEmbedBuilder().
		SetColor(0xED4245).
		SetAuthor(m.Author.Username, "", "").
		SetFields(fields...).
		Build()
	embed.Timestamp.Format(time.RFC3339)

	directMsg := discord.NewMessageCreateBuilder().
		SetEmbeds(embed).
		Build()

	if _, err := b.Client.Rest().CreateMessage(dm.ID(), directMsg); err != nil {
		slog.Error("failed to send bookmark DM", "error", err)
		return err
	}
	return nil
}

func pluralize(s string, n int) string {
	if n == 1 {
		return s
	}
	return s + "s"
}
