package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/txj-xyz/discord-quote-bot/internal/bot"
	"github.com/txj-xyz/discord-quote-bot/internal/config"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create and initialize the bot
	discordBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Load commands
	discordBot.LoadCommands()

	// Start the bot
	if err := discordBot.Start(); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	// Wait for a signal to gracefully shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close the bot connection
	if err := discordBot.Close(); err != nil {
		log.Printf("Error closing bot: %v", err)
	}
} 