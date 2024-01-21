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
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	stg := telestage.NewSceneManager(stateStore, bot)
	mainScene := telestage.NewScene()
	stg.Add("main", mainScene)

	mainScene.Use(addUserBalance)

	mainScene.OnCommand("add", func(ctx telestage.Context) {
		ctx.Reply("working on it....")
	})

	mainScene.OnMessage(func(ctx telestage.Context) {
		account := ctx.Get("account").(*account)
		_, err := ctx.Reply(fmt.Sprintf("Your balance: %d", account.Balance))
		if err != nil {
			log.Println(err)
		}
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds := bot.GetUpdatesChan(u)

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
	return func(ctx telestage.Context) {
		ctx.Set("account", &account{500})
		ef(ctx)
	}
}
