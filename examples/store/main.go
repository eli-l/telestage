package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eli-l/telestage"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
)

func main() {
	stateStore := telestage.NewInMemoryStateStorage()

	cfg := tgbotapi.NewDefaultBotConfig(os.Getenv("BOT_TOKEN"))
	bot := tgbotapi.NewBot(cfg)
	err := bot.Validate()
	if err != nil {
		panic(err)
	}

	stg := telestage.NewSceneManager(stateStore, bot)
	mainScene := telestage.NewScene()
	stg.Add("main", mainScene)

	mainScene.Use(addUserBalance)

	mainScene.OnCommand("add", func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Reply("working on it....")
	})

	mainScene.OnMessage(func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		account := bctx.Get("account").(*account)
		_, err := bctx.Reply(fmt.Sprintf("Your balance: %d", account.Balance))
		if err != nil {
			log.Println(err)
		}
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds, err := tgbotapi.NewPollingHandler(bot, u).InitUpdatesChannel()
	if err != nil {
		panic(err)
	}

	for upd := range upds {
		err := stg.HandleUpdate(upd)
		if err != nil {
			log.Println(err)
		}
	}
}

type account struct {
	Balance int
}

func addUserBalance(ef telestage.EventFn) telestage.EventFn {
	return func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Set("account", &account{500})
		ef(ctx)
	}
}
