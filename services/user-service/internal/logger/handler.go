package logger

import (
	"context"
	"log/slog"
	"os"
)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != "" {
			rec.Add("userID", c.UserID)
		}
		if c.TelegramID != "" {
			rec.Add("telegramID", c.TelegramID)
		}
		if c.RequestID != "" {
			rec.Add("requestID", c.RequestID)
		}
		if c.Extra != nil {
			for k, v := range c.Extra {
				rec.Add(k, v)
			}
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)}
}

func Init(appMode string) {
	var levels = map[string]slog.Leveler{
		"local": slog.LevelDebug,
		"dev":   slog.LevelInfo,
		"prod":  slog.LevelInfo,
	}

	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     levels[appMode],
		AddSource: appMode == "local",
	})

	handler := NewHandlerMiddleware(base).WithAttrs([]slog.Attr{
		slog.String("env", appMode),
	})
	slog.SetDefault(slog.New(handler))
}
