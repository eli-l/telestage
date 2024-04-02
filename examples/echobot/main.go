package main

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/eli-l/telestage"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
)

func main() {
	config := tgbotapi.NewDefaultBotConfig(os.Getenv("BOT_TOKEN"))
	bot := tgbotapi.NewBot(config)

	//bot.Request(tgbotapi.WebhookConfig{
	//	URL: nil,
	//})
	//hook, _ := bot.GetWebhookInfo()
	//_ = hook

	if err := bot.Validate(); err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	rStateStore := telestage.NewRedisStateStorage(rdb, telestage.ExpireHour)
	_ = rStateStore

	stateStore := telestage.NewInMemoryStateStorage()

	mainScene := telestage.NewScene()
	messageScene := telestage.NewScene()

	messageScene.OnStart(func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Reply("Hello from message scene, go back: /leave")
	})

	messageScene.OnCommand("leave", func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		stateStore.Set(ctx, "main")
		bctx.Reply("Welcome in main scene, send: /start")
	})

	messageScene.OnMessage(func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Reply("You in message scene")
	})

	mainScene.OnStart(func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Reply("Hello world. send: /enter")
	})

	mainScene.OnCommand("enter", func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		stateStore.Set(ctx, "message")
		bctx.Reply("Now send: /start")
	})

	mainScene.OnMessage(func(ctx context.Context) {
		bctx := telestage.GetBotContext(ctx)
		bctx.Reply("Incorrect input")
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds, err := tgbotapi.NewPollingHandler(bot, u).InitUpdatesChannel()
	if err != nil {
		log.Fatal(err)
	}

	sceneManager := telestage.NewSceneManagerWithDefault(stateStore, bot, "main")
	sceneManager.Add("main", mainScene)
	sceneManager.Add("message", messageScene)

	for upd := range upds {
		err := sceneManager.HandleUpdate(upd)
		if err != nil {
			log.Println(err)
		}
	}
}
