package log

import (
	"context"
	"os"

	"log/slog"
)

type ContextKeyType int

const ContextKey ContextKeyType = 0

const (
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
)

func New(level slog.Leveler) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     level,
			},
		),
	)
}

func With(logger *slog.Logger, layer, protocol string) *slog.Logger {
	return logger.With("layer", layer).With("protocol", protocol)
}

func FromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(ContextKey).(*slog.Logger)
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ContextKey, logger)
}
