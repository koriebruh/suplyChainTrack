package pkg

import (
	"context"
	"log/slog"
	"os"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (m MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if e := h.Handle(ctx, r); e != nil && err == nil {
				err = e
			}
		}
	}
	return err
}

func (m MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return MultiHandler{handlers: newHandlers}
}

func (m MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return MultiHandler{handlers: newHandlers}
}

func NewSlog() {
	file, err := os.OpenFile("logger/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// Console handler
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// File handler
	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Combine both
	logger := slog.New(MultiHandler{handlers: []slog.Handler{consoleHandler, fileHandler}})
	slog.SetDefault(logger)
	slog.Info("slog internal logging initialized")
}
