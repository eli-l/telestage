package main

import (
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

	mainScene.Use(func(ef telestage.EventFn) telestage.EventFn {
		return func(ctx telestage.Context) {
			if ctx.Message().Sticker == nil { // ignore if message is sticker
				ef(ctx)
			}
		}
	})

	mainScene.OnCommand("ping", func(ctx telestage.Context) {
		ctx.Reply("pong")
	}, func(ef telestage.EventFn) telestage.EventFn {
		return func(ctx telestage.Context) {
			if ctx.Upd().FromChat().IsPrivate() {
				ef(ctx)
			} else {
				ctx.Reply("This command available only in private chat")
			}
		}
	})

	mainScene.OnMessage(func(ctx telestage.Context) {
		ctx.Reply("Hello") // answer on any message
	})

	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds := bot.GetUpdatesChan(u)

	for upd := range upds {
		stg.HandleUpdate(upd)
	}
}
