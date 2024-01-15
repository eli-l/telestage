package telestage

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context interface {
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

// NativeContext ...
type NativeContext struct {
	bot   *tgbotapi.BotAPI
	upd   *tgbotapi.Update
	lock  sync.RWMutex
	store map[string]interface{}

	disableWebPreview bool
}

func (nc *NativeContext) Bot() *tgbotapi.BotAPI {
	return nc.bot
}

func (nc *NativeContext) Upd() *tgbotapi.Update {
	return nc.upd
}

func (nc *NativeContext) Message() *tgbotapi.Message {
	switch {
	case nc.upd.Message != nil:
		return nc.upd.Message
	case nc.upd.CallbackQuery != nil:
		return nc.upd.CallbackQuery.Message
	case nc.upd.EditedMessage != nil:
		return nc.upd.EditedMessage
	case nc.upd.ChannelPost != nil:
		if nc.upd.ChannelPost.PinnedMessage != nil {
			return nc.upd.ChannelPost.PinnedMessage
		}
		return nc.upd.ChannelPost
	case nc.upd.EditedChannelPost != nil:
		return nc.upd.EditedChannelPost
	default:
		return nil
	}
}

func (nc *NativeContext) Sender() *tgbotapi.User {
	switch {
	case nc.upd.CallbackQuery != nil:
		return nc.upd.CallbackQuery.From
	case nc.Message() != nil:
		return nc.Message().From
	case nc.upd.InlineQuery != nil:
		return nc.upd.InlineQuery.From
	case nc.upd.ShippingQuery != nil:
		return nc.upd.ShippingQuery.From
	case nc.upd.PreCheckoutQuery != nil:
		return nc.upd.PreCheckoutQuery.From
	case nc.upd.PollAnswer != nil:
		return nc.upd.PollAnswer.User
	case nc.upd.MyChatMember != nil:
		return &nc.upd.MyChatMember.From
	case nc.upd.ChatMember != nil:
		return &nc.upd.ChatMember.From
	case nc.upd.ChatJoinRequest != nil:
		return &nc.upd.ChatJoinRequest.From
	default:
		return nil
	}
}

func (nc *NativeContext) Chat() *tgbotapi.Chat {
	switch {
	case nc.upd.Message != nil:
		return &nc.upd.Message.Chat
	case nc.Message() != nil:
		return &nc.Message().Chat
	case nc.upd.MyChatMember != nil:
		return &nc.upd.MyChatMember.Chat
	case nc.upd.ChatMember != nil:
		return &nc.upd.ChatMember.Chat
	case nc.upd.ChatJoinRequest != nil:
		return &nc.upd.ChatJoinRequest.Chat
	default:
		return nil
	}
}

func (nc *NativeContext) ChatID() int64 {
	if c := nc.Chat(); c != nil {
		return c.ID
	}

	return nc.Sender().ID
}

func (nc *NativeContext) Text() string {
	m := nc.Message()
	if m == nil {
		return ""
	}
	if m.Caption != "" {
		return m.Caption
	}

	return m.Text
}

func (nc *NativeContext) Reply(text string) (tgbotapi.Message, error) {
	m := tgbotapi.NewMessage(nc.ChatID(), text)
	m.LinkPreviewOptions.IsDisabled = nc.disableWebPreview
	return nc.bot.Send(tgbotapi.NewMessage(nc.ChatID(), text))
}

func (nc *NativeContext) ReplyWithMenu(text string, menu interface{}) (tgbotapi.Message, error) {
	return nc.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: nc.ChatID(),
			},
			ReplyParameters: tgbotapi.ReplyParameters{
				MessageID: 0,
			},
			ReplyMarkup: menu,
		},
		Text: text,
		LinkPreviewOptions: tgbotapi.LinkPreviewOptions{
			IsDisabled: nc.disableWebPreview,
		},
	})
}

func (nc *NativeContext) ReplyHTML(text string) (tgbotapi.Message, error) {
	return nc.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: nc.ChatID(),
			},
			ReplyParameters: tgbotapi.ReplyParameters{
				MessageID: 0,
			},
		},
		ParseMode: tgbotapi.ModeHTML,
		Text:      text,
		LinkPreviewOptions: tgbotapi.LinkPreviewOptions{
			IsDisabled: nc.disableWebPreview,
		},
	})
}

func (nc *NativeContext) ReplyWithMenuHTML(text string, menu interface{}) (tgbotapi.Message, error) {
	return nc.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: nc.ChatID(),
			},
			ReplyParameters: tgbotapi.ReplyParameters{
				MessageID: 0,
			},
			ReplyMarkup: menu,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
		LinkPreviewOptions: tgbotapi.LinkPreviewOptions{
			IsDisabled: nc.disableWebPreview,
		},
	})
}

func (nc *NativeContext) SetDisableWebPreviewForShortMethods(isDisabled bool) {
	nc.disableWebPreview = isDisabled
}

func (nc *NativeContext) Set(key string, value interface{}) {
	nc.lock.Lock()
	defer nc.lock.Unlock()

	if nc.store == nil {
		nc.store = make(map[string]interface{})
	}
	nc.store[key] = value
}

func (nc *NativeContext) Get(key string) interface{} {
	nc.lock.RLock()
	defer nc.lock.RUnlock()
	return nc.store[key]
}
