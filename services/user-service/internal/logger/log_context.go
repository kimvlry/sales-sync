package logger

import "context"

type logCtx struct {
	UserID     string
	TelegramID string
	RequestID  string
	Extra      map[string]any
}

type keyType struct{}

var key = keyType{}

type Builder struct {
	ctx    context.Context
	logCtx logCtx
}

func NewBuilder(ctx context.Context) *Builder {
	existing := logCtx{}
	if c, ok := ctx.Value(key).(logCtx); ok {
		existing = c
	}
	return &Builder{
		ctx:    ctx,
		logCtx: existing,
	}
}

func (b *Builder) WithUserID(id string) *Builder {
	b.logCtx.UserID = id
	return b
}

func (b *Builder) WithTelegramID(tid string) *Builder {
	b.logCtx.TelegramID = tid
	return b
}

func (b *Builder) WithRequestID(reqID string) *Builder {
	b.logCtx.RequestID = reqID
	return b
}

func (b *Builder) WithExtra(extra map[string]any) *Builder {
	if b.logCtx.Extra == nil {
		b.logCtx.Extra = make(map[string]any)
	}
	for k, v := range extra {
		b.logCtx.Extra[k] = v
	}
	return b
}

func (b *Builder) Build() context.Context {
	return context.WithValue(b.ctx, key, b.logCtx)
}
