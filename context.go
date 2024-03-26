package telestage

import (
	"context"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
)

type BotCtxTyp string

const BotCtxKey BotCtxTyp = "botCtx"

type BotContext interface {
	// Bot ...
	Bot() *tgbotapi.BotAPI
	// Upd ...
	Upd() *tgbotapi.Update
	// Message ...
	Message() *tgbotapi.Message
	// Sender ...
	Sender() *tgbotapi.User
	// Chat ...
	Chat() *tgbotapi.Chat
	// ChatID ...
	ChatID() int64
	// Text ...
	Text() string

	// Fast methods
	// SetDisableWebPreviewForShortMethods ...
	SetDisableWebPreviewForShortMethods(bool)
	// Reply ...
	Reply(string) (tgbotapi.Message, error)
	// ReplyWithMenu ...
	ReplyWithMenu(string, interface{}) (tgbotapi.Message, error)
	// ReplyHTML ...
	ReplyHTML(string) (tgbotapi.Message, error)
	// ReplyWithMenuHTML ...
	ReplyWithMenuHTML(string, interface{}) (tgbotapi.Message, error)

	// Get retrieves data from the context.
	Get(key string) interface{}
	// Set saves data in the context.
	Set(key string, val interface{})
}

func WithBotContext(ctx context.Context, botCtx BotContext) context.Context {
	return context.WithValue(ctx, BotCtxKey, botCtx)
}

func GetBotContext(ctx context.Context) BotContext {
	return ctx.Value(BotCtxKey).(BotContext)
}
