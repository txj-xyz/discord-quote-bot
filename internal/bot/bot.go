package bot

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/txj-xyz/discord-quote-bot/internal/commands"
	"github.com/txj-xyz/discord-quote-bot/internal/config"
)

// Command represents a slash command with its handler
type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Bot represents the Discord bot instance
type Bot struct {
	Session  *discordgo.Session
	Config   *config.Config
	Commands map[string]*Command
	mu       sync.RWMutex
}

// New creates a new bot instance
func New(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	bot := &Bot{
		Session:  session,
		Config:   cfg,
		Commands: make(map[string]*Command),
	}

	// Add handlers
	session.AddHandler(bot.handleInteractionCreate)
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot is up as %s#%s", r.User.Username, r.User.Discriminator)
	})

	// Set required intents
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	return bot, nil
}

// RegisterCommand registers a new command with the bot
func (b *Bot) RegisterCommand(cmd *Command) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Commands[cmd.ApplicationCommand.Name] = cmd
}

// Start starts the bot
func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	// Register commands with Discord API
	for _, cmd := range b.Commands {
		var err error
		if b.Config.Bot.GuildID != "" {
			// Guild command for testing
			_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.Config.Bot.GuildID, cmd.ApplicationCommand)
		} else {
			// Global command for production
			_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd.ApplicationCommand)
		}

		if err != nil {
			return fmt.Errorf("error creating command '%s': %w", cmd.ApplicationCommand.Name, err)
		}
		log.Printf("Registered command: %s", cmd.ApplicationCommand.Name)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	return nil
}

// Close closes the bot connection
func (b *Bot) Close() error {
	return b.Session.Close()
}

// LoadCommands loads all available commands
func (b *Bot) LoadCommands() {
	// Register the quote command
	quoteCmd := &Command{
		ApplicationCommand: commands.NewQuoteCommand(b.Config.Quote.Channel),
		Handler:            commands.HandleQuoteCommand(b.Config.Quote.Channel),
	}
	b.RegisterCommand(quoteCmd)
}

// handleInteractionCreate handles interaction create events
func (b *Bot) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only handle application commands
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	b.mu.RLock()
	cmd, exists := b.Commands[i.ApplicationCommandData().Name]
	b.mu.RUnlock()

	if !exists {
		log.Printf("Unknown command: %s", i.ApplicationCommandData().Name)
		return
	}

	cmd.Handler(s, i)
} 