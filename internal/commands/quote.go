package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// NewQuoteCommand creates a new quote command
func NewQuoteCommand(quoteChannelID string) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "quote",
		Description: "Quote a user with a message that will be posted to the quotes channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to quote",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The quote message",
				Required:    true,
			},
		},
	}
}

// HandleQuoteCommand handles the quote command
func HandleQuoteCommand(quoteChannelID string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Acknowledge the interaction immediately
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Error responding to interaction: %v", err)
			return
		}

		// Get options from the command
		options := i.ApplicationCommandData().Options
		var userID, quote string

		for _, option := range options {
			switch option.Name {
			case "user":
				userID = option.Value.(string)
			case "message":
				quote = option.Value.(string)
			}
		}

		// Create the embed
		embed := &discordgo.MessageEmbed{
			Title:       "Quote",
			Description: quote,
			Color:       0x00AAFF, // Blue color
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Quoted by %s", i.Member.User.Username),
			},
		}

		// Send the message to the quote channel
		_, err = s.ChannelMessageSendComplex(quoteChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s> was quoted:", userID),
			Embeds:  []*discordgo.MessageEmbed{embed},
		})

		if err != nil {
			log.Printf("Error sending quote message: %v", err)
			// Respond with an error message
			errorMsg := fmt.Sprintf("Error sending quote: %v", err)
			_, respErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &errorMsg,
			})
			if respErr != nil {
				log.Printf("Error editing interaction response: %v", respErr)
			}
			return
		}

		// Send a success message to the user
		successMsg := "Quote posted successfully!"
		_, respErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &successMsg,
		})
		if respErr != nil {
			log.Printf("Error editing interaction response: %v", respErr)
		}
	}
} 