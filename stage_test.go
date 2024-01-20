package telestage

//
//import (
//	"testing"
//
//	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
//	"github.com/stretchr/testify/assert"
//)
//
//type stateStorage struct {
//	state string
//}
//
//func (s *stateStorage) Set(state string) {
//	s.state = state
//}
//
//func (s *stateStorage) Get() string {
//	return s.state
//}
//func TestStateGetting(t *testing.T) {
//	firstScene := NewScene()
//	firstSceneInvoked := false
//	firstScene.OnMessage(func(ctx Context) {
//		firstSceneInvoked = true
//	})
//
//	secondScene := NewScene()
//	secondSceneInvoked := false
//	secondScene.OnMessage(func(ctx Context) {
//		secondSceneInvoked = true
//	})
//	stateStorage := &stateStorage{}
//	stateGetter := func(ctx Context) string {
//		return stateStorage.Get()
//	}
//	stage := NewSceneManager(stateGetter)
//
//	stage.Add("", firstScene)
//	stage.Add("second", secondScene)
//
//	stateStorage.Set("")
//	stage.HandleUpdate(&tgbotapi.BotAPI{}, tgbotapi.Update{
//		Message: &tgbotapi.Message{},
//	})
//	assert.Condition(t, func() (success bool) {
//		if firstSceneInvoked && !secondSceneInvoked {
//			return true
//		}
//		return false
//	}, "if state = '', only firstScene event(OnMessage) must be invoked")
//
//	stateStorage.Set("second")
//	firstSceneInvoked = false
//	secondSceneInvoked = false
//	stage.HandleUpdate(&tgbotapi.BotAPI{}, tgbotapi.Update{
//		Message: &tgbotapi.Message{},
//	})
//
//	assert.Condition(t, func() (success bool) {
//		if !firstSceneInvoked && secondSceneInvoked {
//			return true
//		}
//		return false
//	}, "if state = 'second', only secondScene event(OnMessage) must be invoked", firstSceneInvoked, secondSceneInvoked)
//}
//
//func TestUndefinedScene(t *testing.T) {
//	stage := NewSceneManager(emptyStateGetter)
//
//	err := stage.HandleUpdate(&tgbotapi.BotAPI{}, tgbotapi.Update{})
//	assert.Error(t, err, "call undefined scene")
//
//}
//
//func TestStage_Add(t *testing.T) {
//	stage := NewSceneManager(emptyStateGetter)
//
//	scene := NewScene()
//	k := "main"
//	stage.Add(k, scene)
//
//	assert.Equal(t, scene, stage.scenes[k], "add scene to stage")
//}
//
//func TestStage_NewStage(t *testing.T) {
//	sg := func(ctx Context) string {
//		return "test"
//	}
//	stage := NewSceneManager(sg)
//	ctx := &NativeContext{
//		bot: &tgbotapi.BotAPI{},
//		upd: &tgbotapi.Update{},
//	}
//	assert.Equal(t, sg(ctx), stage.stateGetter(ctx), "new stage")
//}
