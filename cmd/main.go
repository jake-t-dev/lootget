package main

import (
	"context"
	"log"

	"github.com/jake-t-dev/lootget/internal/bot"
	"github.com/jake-t-dev/lootget/internal/config"
)

func main() {
	ctx := context.Background()

	ctx, err := config.NewConfig(ctx)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	token := ctx.Value("config").(*config.Config).Token
	b := &bot.Bot{
		Token: token,
	}

	b.Start()
}
