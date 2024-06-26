package telestage

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
)

var (
	ErrSceneNotFound = errors.New("scene not found")
)

type ErrorHandler func(ctx context.Context, err error) error

func DefaultErrorHandler(ctx context.Context, err error) error {
	fmt.Println(err)
	return err
}

type SceneManagerInterface interface {
	Add(state string, scene *Scene)
	HandleUpdate(upd tgbotapi.Update) error
}

type SceneManager struct {
	scenes       map[State]*Scene
	bot          *tgbotapi.BotAPI
	stateStorage StateStorage
	defState     State
	ErrorHandler ErrorHandler
}

func NewSceneManager(storage StateStorage, bot *tgbotapi.BotAPI) *SceneManager {
	return &SceneManager{
		scenes:       map[State]*Scene{"": {}},
		bot:          bot,
		stateStorage: storage,
		defState:     "",
		ErrorHandler: DefaultErrorHandler,
	}
}

func NewSceneManagerWithDefault(storage StateStorage, bot *tgbotapi.BotAPI, defState string) *SceneManager {
	return &SceneManager{
		scenes:       map[State]*Scene{},
		bot:          bot,
		stateStorage: storage,
		defState:     State(defState),
		ErrorHandler: DefaultErrorHandler,
	}
}

func (s *SceneManager) Add(state string, scene *Scene) {
	st := State(state)
	s.scenes[st] = scene
}

func (s *SceneManager) Get(sc string) *Scene {
	st := State(sc)
	scene, ok := s.scenes[st]
	if !ok {
		return &Scene{}
	}
	return scene
}

func (s *SceneManager) SetErrorHandler(handler ErrorHandler) {
	s.ErrorHandler = handler
}

func (s *SceneManager) HandleUpdate(upd tgbotapi.Update) error {
	ctx := context.WithValue(context.Background(), BotCtxKey, &NativeContext{
		bot: s.bot,
		upd: &upd,
	})

	state, err := s.stateStorage.Get(ctx)
	if err != nil {
		return err
	}

	var scene *Scene
	var ok bool
	if state == "" {
		scene, ok = s.scenes[s.defState]
	} else {
		scene, ok = s.scenes[state]
	}

	if !ok {
		return s.ErrorHandler(ctx, fmt.Errorf("%w with name %s", ErrSceneNotFound, state))
	}

	events := scene.GetEventHandler()
	for _, e := range events {
		if e(ctx) {
			return nil
		}
	}

	return nil
}
