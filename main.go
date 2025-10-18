package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env: %v\n", err)
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error starting up discord client: %v\n", err)
	}

	dg.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentDirectMessages
	dg.AddHandlerOnce(ready)
	dg.AddHandler(handleInteraction)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord Bot connection: %v\n", err)
	}
	defer dg.Close()
	cleanupCommands(dg, nil)

	rcmdID, err := registerCommands(dg)
	if err != nil {
		log.Println(err)
		return
	}

	defer cleanupCommands(dg, rcmdID)

	fmt.Println("Bot running. Press CTRL+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("Shutting down bot")
}

func cleanupCommands(s *discordgo.Session, c *discordgo.ApplicationCommand) error {
	appID := os.Getenv("APPLICATION_ID")
	guildID := os.Getenv("GUILD_ID")
	if c != nil {
		err := s.ApplicationCommandDelete(appID, guildID, c.ID)
		if err != nil {
			return err
		}
	}
	cmds, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		return err
	}
	for _, v := range cmds {
		err := s.ApplicationCommandDelete(v.ApplicationID, v.GuildID, v.ID)
		if err != nil {
			return err
		}

	}
	return err
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
}

func registerCommands(s *discordgo.Session) (*discordgo.ApplicationCommand, error) {
	cmdName := "Bookmark local"
	command := &discordgo.ApplicationCommand{
		Name: cmdName,
		Type: discordgo.MessageApplicationCommand,
	}
	appID := os.Getenv("APPLICATION_ID")
	guildID := os.Getenv("GUILD_ID")
	rcmdID, err := s.ApplicationCommandCreate(appID, guildID, command)

	if err != nil {
		log.Printf("Error while creating the Message Context Menu command. Command name: %v, err: %v\n", cmdName, err)
		return nil, err
	}

	return rcmdID, nil
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()
	msg := data.Resolved.Messages[data.TargetID]
	if msg == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not find the message.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	dmChannel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("Failed to create DM channel: %v", err)
		return
	}
	msgChannel, err := s.Channel(msg.ChannelID)
	if err != nil {
		fmt.Printf("Failed to fetch channel: %v", err)
		return
	}
	msgGuild, err := s.Guild(msgChannel.GuildID)
	if err != nil {
		fmt.Printf("Failed to fetch guild: %v", err)
		return
	}
	responseMsg := &discordgo.MessageEmbed{
		// Title:       "You bookmarked a message",
		// Description: fmt.Sprintf("> %s", msg.Content),
		Color: 0xED4245,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.Username,
			IconURL: msg.Author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "From",
				Value: fmt.Sprintf("%s, #%s", msgGuild.Name, msgChannel.Name),
			},
			{
				Value: fmt.Sprintf("[Link to the original message](https://discord.com/channels/%s/%s/%s)", msgGuild.ID, msg.ChannelID, msg.ID),
			},
			{
				Name:  "Message preview",
				Value: msg.Content,
			},
			func() *discordgo.MessageEmbedField {
				n := len(msg.Attachments)
				if n == 0 {
					return &discordgo.MessageEmbedField{}
				}

				return &discordgo.MessageEmbedField{
					Value: fmt.Sprintf("%d %s", n, pluralize("attachment", n)),
				}
			}(),
			func() *discordgo.MessageEmbedField {
				n := len(msg.Embeds)
				if n == 0 {
					return &discordgo.MessageEmbedField{}
				}

				return &discordgo.MessageEmbedField{
					Value: fmt.Sprintf("%d %s", n, pluralize("embed", n)),
				}
			}(),
			// {
			// 	Value: fmt.Sprintf("%d %s", len(msg.Embeds), pluralize("attachment", len(msg.Attachments))),
			// },
			// {
			// 	Value: fmt.Sprintf("%d %s", len(msg.Embeds), pluralize("embed", len(msg.Embeds))),
			// },
		},
		Timestamp: msg.Timestamp.Format(time.RFC3339),
	}
	_, err = s.ChannelMessageSendEmbed(dmChannel.ID, responseMsg)
	if err != nil {
		log.Printf("Failed to send bookmark DM: %v", err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Bookmark sent to your DMs",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

}

func pluralize(s string, n int) string {
	if n == 1 {
		return s
	}

	return s + "s"
}
