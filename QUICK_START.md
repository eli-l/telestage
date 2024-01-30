# Quick start

### Scene middlewares

```go
mainScene.Use(func(ef telestage.EventFn) telestage.EventFn {
	return func(ctx telestage.Context) {
		if ctx.Message().Sticker == nil { // ignore if message is sticker
			ef(ctx)
		}
	}
})

mainScene.OnMessage(func(ctx telestage.Context) {
    ctx.Reply("Hello") // answer on any message
})
```

### Event group middlewares

```go
mainScene.UseGroup(func(s *telestage.Scene) {
    // s is mainScene
    s.OnCommand("ban", func(ctx telestage.Context) {...})
    s.OnCommand("kick", func(ctx telestage.Context) {...})
}, func(ef telestage.EventFn) telestage.EventFn {
    return func(ctx telestage.Context) {
        if isAdmin(ctx.Sender().ID) {
            ef(ctx)
        }
    }
})
```

### Event middlewares

```go
mainScene.OnCommand("ping", func(ctx telestage.Context) {
    ctx.Reply("pong")
}, func(ef telestage.EventFn) telestage.EventFn {
    return func(ctx telestage.Context) {
        if ctx.Upd().FromChat().IsPrivate() {
            ef(ctx)
        } else {
            ctx.Reply("This command available only in private chat")
        }
    }
})
```
