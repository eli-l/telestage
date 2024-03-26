package telestage

import "context"

type EventFn func(context.Context)
type EventHandler func(context.Context) bool
type EventDeterminant func(context.Context) bool

type Scene struct {
	eventHandlers []EventHandler
	middlewares   []Middleware
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) GetEventHandler() []EventHandler {
	return s.eventHandlers
}

func (s *Scene) Use(mw ...Middleware) {
	s.middlewares = append(s.middlewares, mw...)
}

func (s *Scene) UseGroup(group func(*Scene), mw ...Middleware) {
	original := s.middlewares
	s.middlewares = append(s.middlewares, mw...)
	group(s)
	s.middlewares = original
}

// OnCommand handle the command specified by first argument
func (s *Scene) OnCommand(cmd string, ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.eventHandlers = append(s.eventHandlers, func(ctx context.Context) bool {
		bctx := GetBotContext(ctx)
		if bctx.Upd().Message == nil || bctx.Upd().Message.Command() != cmd {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnMessage handle any message type (photo, text, sticker etc.)
func (s *Scene) OnMessage(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.eventHandlers = append(s.eventHandlers, func(ctx context.Context) bool {
		bctx := GetBotContext(ctx)
		if bctx.Upd().Message == nil {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnPhoto handle sending a photo
func (s *Scene) OnPhoto(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.eventHandlers = append(s.eventHandlers, func(ctx context.Context) bool {
		bctx := GetBotContext(ctx)
		m := bctx.Message()
		if m == nil || len(m.Photo) == 0 {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnSticker handle sending a sticker
func (s *Scene) OnSticker(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.eventHandlers = append(s.eventHandlers, func(ctx context.Context) bool {
		bctx := GetBotContext(ctx)
		m := bctx.Message()
		if m == nil || m.Sticker == nil {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnPhoto handle the "/start" command
func (s *Scene) OnStart(ef EventFn, mw ...Middleware) {
	s.OnCommand("start", ef, mw...)
}

// On handle the your own event determinator
func (s *Scene) On(determinant EventDeterminant, ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.eventHandlers = append(s.eventHandlers, func(ctx context.Context) bool {
		if !determinant(ctx) {
			return false
		}
		ef(ctx)

		return true
	})
}
