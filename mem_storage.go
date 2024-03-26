package telestage

import "context"

type InMemoryStateStorage struct {
	states map[int64]State
}

func NewInMemoryStateStorage() *InMemoryStateStorage {
	return &InMemoryStateStorage{
		states: map[int64]State{},
	}
}

func (m *InMemoryStateStorage) Get(ctx context.Context) (State, error) {
	bctx := GetBotContext(ctx)

	if bctx.Sender() == nil {
		return "", ErrNoSender
	}

	state, ok := m.states[bctx.Sender().ID]
	if !ok {
		return "main", nil
	}

	return state, nil
}

func (m *InMemoryStateStorage) Set(ctx context.Context, state State) error {
	bctx := GetBotContext(ctx)
	m.states[bctx.Sender().ID] = state
	return nil
}
