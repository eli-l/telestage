package telestage

import "errors"

var ErrNoSender = errors.New("no sender available in Update")

type State string

func (s *State) String() string {
	return string(*s)
}

type StateStorage interface {
	Get(ctx Context) (State, error)
	Set(ctx Context, state State) error
}

type InMemoryStateStorage struct {
	states map[int64]State
}

func NewInMemoryStateStorage() *InMemoryStateStorage {
	return &InMemoryStateStorage{
		states: map[int64]State{},
	}
}

func (m *InMemoryStateStorage) Get(ctx Context) (State, error) {
	if ctx.Sender() == nil {
		return "", ErrNoSender
	}

	state, ok := m.states[ctx.Sender().ID]
	if !ok {
		return "main", nil
	}

	return state, nil
}

func (m *InMemoryStateStorage) Set(ctx Context, state State) error {
	m.states[ctx.Sender().ID] = state
	return nil
}
