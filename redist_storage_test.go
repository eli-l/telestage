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

package telestage_test

import (
	"context"
	"testing"

	tgbotapi "github.com/eli-l/telegram-bot-api/v7"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	_ "github.com/go-redis/redismock/v9"

	"github.com/eli-l/telestage"
)

func Test_RedisStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	bctx := NewMockBotContext(ctrl)
	bctx.EXPECT().Sender().Return(&tgbotapi.User{ID: 1}).AnyTimes()
	ctx = telestage.WithBotContext(ctx, bctx)
	db, mock := redismock.NewClientMock()

	t.Run("Set", func(t *testing.T) {
		mock.ExpectSet("1", "newState", -1).SetVal("OK")
		storage := telestage.NewRedisStateStorage(db, telestage.ExpireNever)

		err := storage.Set(ctx, "newState")
		require.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		mock.ExpectGet("1").SetVal("newState")
		storage := telestage.NewRedisStateStorage(db, telestage.ExpireNever)

		state, err := storage.Get(ctx)
		require.NoError(t, err)
		require.Equal(t, telestage.State("newState"), state)
	})
}
