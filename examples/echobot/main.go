package main

import (
	"log"
	"os"

	"github.com/eli-l/telestage"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	stateStore := telestage.NewInMemoryStateStorage()

	mainScene := telestage.NewScene()
	messageScene := telestage.NewScene()

	messageScene.OnStart(func(ctx telestage.Context) {
		ctx.Reply("Hello from message scene, go back: /leave")
	})

	messageScene.OnCommand("leave", func(ctx telestage.Context) {
		stateStore.Set(ctx, "main")
		ctx.Reply("Welcome in main scene, send: /start")
	})

	messageScene.OnMessage(func(ctx telestage.Context) {
		ctx.Reply("You in message scene")
	})

	mainScene.OnStart(func(ctx telestage.Context) {
		ctx.Reply("Hello world. send: /enter")
	})

	mainScene.OnCommand("enter", func(ctx telestage.Context) {
		stateStore.Set(ctx, "message")
		ctx.Reply("Now send: /start")
	})

	mainScene.OnMessage(func(ctx telestage.Context) {
		ctx.Reply("Incorrect input")
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds := bot.GetUpdatesChan(u)

	sceneManager := telestage.NewSceneManager(stateStore, bot)
	sceneManager.Add("main", mainScene)
	sceneManager.Add("message", messageScene)

	for upd := range upds {
		err := sceneManager.HandleUpdate(upd)
		if err != nil {
			log.Println(err)
		}
	}
}
