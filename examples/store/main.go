package main

import (
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

	mainScene.OnCommand("add", func(ctx telestage.BotContext) {
		ctx.Reply("working on it....")
	})

	mainScene.OnMessage(func(ctx telestage.BotContext) {
		account := ctx.Get("account").(*account)
		_, err := ctx.Reply(fmt.Sprintf("Your balance: %d", account.Balance))
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
	return func(ctx telestage.BotContext) {
		ctx.Set("account", &account{500})
		ef(ctx)
	}
}
