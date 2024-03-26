package telestage

import (
	"context"
	"errors"
)

var ErrNoSender = errors.New("no sender available in Update")

type State string

func (s *State) String() string {
	return string(*s)
}

type StateStorage interface {
	Get(ctx context.Context) (State, error)
	Set(ctx context.Context, state State) error
}
