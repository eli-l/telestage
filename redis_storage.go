// MIT License
//
// Copyright (c) 2024 eli-l
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package telestage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ExpireExpired = iota
	ExpireNever   = iota
	ExpireKeep    = iota
	ExpireHour    = iota
	ExpireDay     = iota
)

func expirationFromInt(expiration int) time.Duration {
	switch expiration {
	case ExpireExpired:
		return 0
	case ExpireHour:
		return time.Hour
	case ExpireDay:
		return time.Hour * 24
	case ExpireNever, ExpireKeep:
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
	err := s.client.Set(ctx, id, state.String(), s.defaultExpiration).Err()
	return err
}
