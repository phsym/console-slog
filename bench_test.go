package console

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
)

type DummyHandler struct{}

func (*DummyHandler) Enabled(context.Context, slog.Level) bool   { return true }
func (*DummyHandler) Handle(context.Context, slog.Record) error  { return nil }
func (h *DummyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *DummyHandler) WithGroup(name string) slog.Handler       { return h }

var handlers = []struct {
	name string
	hdl  slog.Handler
}{
	{"dummy", &DummyHandler{}},
	{"console", NewHandler(io.Discard, &HandlerOptions{Level: slog.LevelDebug, AddSource: false})},
	{"std-text", slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})},
	{"std-json", slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})},
}

var attrs = []slog.Attr{
	slog.String("titi", "toto"),
	slog.String("tata", "tutu"),
	slog.Int("foo", 12),
	slog.Duration("dur", 3*time.Second),
	slog.Bool("bool", true),
	slog.Float64("float", 23.7),
	slog.Time("thetime", time.Now()),
	slog.Any("err", errors.New("yo")),
	slog.Group("empty"),
	slog.Group("group", slog.String("bar", "baz")),
}

var attrsAny = []any{
	slog.String("titi", "toto"),
	slog.String("tata", "tutu"),
	slog.Int("foo", 12),
	slog.Duration("dur", 3*time.Second),
	slog.Bool("bool", true),
	slog.Float64("float", 23.7),
	slog.Time("thetime", time.Now()),
	slog.Any("err", errors.New("yo")),
	slog.Group("empty"),
	slog.Group("group", slog.String("bar", "baz")),
}

func BenchmarkHandlers(b *testing.B) {
	ctx := context.Background()
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "yo", 0)
	rec.AddAttrs(attrs...)

	for _, tc := range handlers {
		b.Run(tc.name, func(b *testing.B) {
			// l := tc.hdl.WithAttrs([]slog.Attr{slog.String("toto", "tutu")}) //.WithGroup("test")
			l := tc.hdl.WithAttrs(attrs).WithGroup("test").WithAttrs(attrs)
			// Warm-up
			l.Handle(ctx, rec)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				l.Handle(ctx, rec)
			}
		})
	}
}

func BenchmarkLoggers(b *testing.B) {
	for _, tc := range handlers {
		ctx := context.Background()
		b.Run(tc.name, func(b *testing.B) {
			l := slog.New(tc.hdl).With(attrsAny...).WithGroup("test").With(attrsAny...)
			// Warm-up
			l.LogAttrs(ctx, slog.LevelInfo, "yo", attrs...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				l.LogAttrs(ctx, slog.LevelInfo, "yo", attrs...)
			}
		})
	}
}
