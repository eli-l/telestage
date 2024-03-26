package telestage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	EXPIRE_EXPIRED = iota
	EXPIRE_NEVER   = iota
	EXPIRE_KEEP    = iota
	EXPIRE_HOUR    = iota
	EXPIRE_DAY     = iota
)

func expirationFromInt(expiration int) time.Duration {
	switch expiration {
	case EXPIRE_EXPIRED:
		return 0
	case EXPIRE_HOUR:
		return time.Hour
	case EXPIRE_DAY:
		return time.Hour * 24
	case EXPIRE_NEVER, EXPIRE_KEEP:
		return -1
	default:
		return -1
	}
}

type RedisStateStorage struct {
	client            *redis.Client
	defaultExpiration time.Duration
}

func NewRedisStateStorage(client *redis.Client, expiration int) *RedisStateStorage {
	return &RedisStateStorage{
		client:            client,
		defaultExpiration: expirationFromInt(expiration),
	}
}

func (s *RedisStateStorage) Get(ctx context.Context) (State, error) {
	botCxt := GetBotContext(ctx)
	id := strconv.Itoa(int(botCxt.Sender().ID))
	r := s.client.Get(ctx, id)
	if r.Err() != nil {
		return "", r.Err()
	}
	return State(r.Val()), nil
}

func (s *RedisStateStorage) Set(ctx context.Context, state State) error {
	botCtx := GetBotContext(ctx)
	id := strconv.Itoa(int(botCtx.Sender().ID))
	err := s.client.Set(ctx, id, state, s.defaultExpiration).Err()
	return err
}
