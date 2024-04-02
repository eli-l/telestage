package telestage_test

import (
	"context"
	"testing"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/eli-l/telestage"
)

func TestStateGetting(t *testing.T) {

	ctrl := gomock.NewController(t)
	ctx := context.Background()
	bctx := NewMockBotContext(ctrl)
	bctx.EXPECT().Sender().Return(&tgbotapi.User{ID: 1}).AnyTimes()
	ctx = telestage.WithBotContext(ctx, bctx)

	firstScene := telestage.NewScene()

	firstSceneInvoked := false
	firstScene.OnMessage(func(ctx context.Context) {
		firstSceneInvoked = true
	})

	secondScene := telestage.NewScene()
	secondSceneInvoked := false
	secondScene.OnMessage(func(ctx context.Context) {
		secondSceneInvoked = true
	})

	storage := telestage.NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}

	sm := telestage.NewSceneManager(storage, bot)

	sm.Add("main", firstScene)
	sm.Add("second", secondScene)

	err := storage.Set(ctx, "main")
	require.NoError(t, err)

	err = sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})
	require.NoError(t, err)

	err = storage.Set(ctx, "second")
	require.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		if firstSceneInvoked && !secondSceneInvoked {
			return true
		}
		return false
	}, "if state = '', only firstScene event(OnMessage) must be invoked")

	firstSceneInvoked = false
	secondSceneInvoked = false
	err = sm.HandleUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				ID: 1,
			},
		},
	})

	require.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		if !firstSceneInvoked && secondSceneInvoked {
			return true
		}
		return false
	}, "if state = 'second', only secondScene event(OnMessage) must be invoked", firstSceneInvoked, secondSceneInvoked)
}

func TestUndefinedScene(t *testing.T) {
	storage := telestage.NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := telestage.NewSceneManager(storage, bot)
	err := sm.HandleUpdate(tgbotapi.Update{})
	require.Error(t, err)
}

func TestStage_Add(t *testing.T) {
	storage := telestage.NewInMemoryStateStorage()
	bot := &tgbotapi.BotAPI{}
	sm := telestage.NewSceneManager(storage, bot)
	scene := telestage.NewScene()
	k := "main"
	sm.Add(k, scene)
	require.Equal(t, scene, sm.Get(k))
}
