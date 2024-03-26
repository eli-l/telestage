package telestage

import (
	"context"
	"strings"
	"testing"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	s := NewScene()
	s.OnMessage(func(_ context.Context) {})
	s.OnStart(func(_ context.Context) {})

	assert.Equal(t, len(s.GetEventHandler()), 2, "get eventHandlers len should be equal with 2")
}

func TestOnCommand(t *testing.T) {
	s := NewScene()
	invoked := false
	cmd := "test"
	s.OnCommand(cmd, func(_ context.Context) {
		invoked = true
	})

	storage := NewInMemoryStateStorage()

	stage := NewSceneManager(storage, &tgbotapi.BotAPI{})
	stage.Add("main", s)

	err := stage.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/" + cmd,
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(cmd) + 1, // plus slash
				},
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, invoked)
}

func TestSceneMiddleware(t *testing.T) {
	s := NewScene()
	s.Use(func(ef EventFn) EventFn {
		return func(ctx context.Context) {
			bctx := GetBotContext(ctx)
			if bctx.Upd().FromChat().IsPrivate() {
				ef(ctx)
			}
		}
	})
	invoked := false
	s.OnMessage(func(_ context.Context) {
		invoked = true
	})
	storage := NewInMemoryStateStorage()

	stage := NewSceneManager(storage, &tgbotapi.BotAPI{})
	stage.Add("main", s)

	err := stage.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: tgbotapi.Chat{
				Type: "group",
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})

	require.NoError(t, err)

	assert.False(t, invoked)

	err = stage.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: tgbotapi.Chat{
				Type: "private",
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)

	assert.True(t, invoked)
}

func TestEventGroupMiddleware(t *testing.T) {
	s := NewScene()

	groupMiddlewareInvoked := false
	s.UseGroup(func(s *Scene) {
		s.OnSticker(func(_ context.Context) {})
	}, func(ef EventFn) EventFn {
		return func(ctx context.Context) {
			groupMiddlewareInvoked = true
			ef(ctx)
		}
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}

	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)

	assert.False(t, groupMiddlewareInvoked, "they should be false, because event OnSticker not invoked")

	err = sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Sticker: &tgbotapi.Sticker{},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)

	assert.True(t, groupMiddlewareInvoked, "they should be true, because event OnSticker invoked")
}

func TestEventMiddleware(t *testing.T) {
	s := NewScene()

	eventMiddlewareInvoked := false
	s.OnMessage(func(_ context.Context) {}, func(ef EventFn) EventFn {
		return func(ctx context.Context) {
			eventMiddlewareInvoked = true
		}
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, eventMiddlewareInvoked, "they should be true, because event OnMessage invoked")
}

func TestOnStart(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnStart(func(_ context.Context) {
		invoked = true
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6, // /start
				},
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message command is /start")
}

func TestOnPhoto(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnPhoto(func(_ context.Context) {
		invoked = true
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Photo: []tgbotapi.PhotoSize{
				{
					FileID: "random_file_id",
				},
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message photo is not nil")
}

func TestOnSticker(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnSticker(func(_ context.Context) {
		invoked = true
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Sticker: &tgbotapi.Sticker{
				FileID: "random_file_id",
			},
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message sticker is not nil")
}

func TestOnMessage(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnMessage(func(ctx context.Context) {
		invoked = true
	})

	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message is not nil")
}

func TestOwnEvent(t *testing.T) {
	s := NewScene()
	messageTextContains := func(text string) EventDeterminant {
		return func(ctx context.Context) bool {
			bctx := GetBotContext(ctx)
			return strings.Contains(bctx.Text(), text)
		}
	}

	invoked := false
	s.On(messageTextContains("hello"), func(_ context.Context) {
		invoked = true
	})
	storage := NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := NewSceneManager(storage, bot)
	sm.Add("main", s)

	err := sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "hello, my name is John Doe",
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message text contains 'hello'")

	invoked = false
	err = sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Caption: "hello, its my first picture",
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, invoked, "they should be true if message caption contains 'hello'")
}
