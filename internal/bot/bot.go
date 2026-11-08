package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jake-t-dev/lootget/internal/bot/handlers"
)

type Bot struct {
	Token string
}

func (b *Bot) init() (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	dg, err = b.addHandlers(dg)
	if err != nil {
		return nil, fmt.Errorf("error adding handlers: %w", err)
	}

	return dg, nil
}

func (b *Bot) addHandlers(dg *discordgo.Session) (*discordgo.Session, error) {
	return handlers.AddHandlers(dg)
}

func (b *Bot) Start() {
	dg, err := b.init()
	if err != nil {
		log.Fatalf("Bot initialization failed: %v", err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
	fmt.Println("Bot gracefully shut down.")
}
