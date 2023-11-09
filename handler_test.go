package console

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestHandler_colors(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, nil)
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", 0)
	AssertNoError(t, h.Handle(context.Background(), rec))

	expected := fmt.Sprintf("\x1b[90m%s\x1b[0m \x1b[92mINF\x1b[0m \x1b[97mfoobar\x1b[0m\r\n", now.Format(time.DateTime))
	AssertEqual(t, expected, buf.String())
}

func TestHandler_NoColor(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, &HandlerOptions{NoColor: true})
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", 0)
	AssertNoError(t, h.Handle(context.Background(), rec))

	expected := fmt.Sprintf("%s INF foobar\r\n", now.Format(time.DateTime))
	AssertEqual(t, expected, buf.String())
}

type theStringer struct{}

func (t theStringer) String() string { return "stringer" }

type noStringer struct {
	Foo string
}

func TestHandler_Attr(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, &HandlerOptions{NoColor: true})
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", 0)
	rec.AddAttrs(
		slog.Bool("bool", true),
		slog.Int("int", -12),
		slog.Uint64("uint", 12),
		slog.Float64("float", 3.14),
		slog.String("foo", "bar"),
		slog.Time("time", now),
		slog.Duration("dur", time.Second),
		slog.Group("group", slog.String("foo", "bar"), slog.Group("subgroup", slog.String("foo", "bar"))),
		slog.Any("err", errors.New("the error")),
		slog.Any("stringer", theStringer{}),
		slog.Any("nostringer", noStringer{Foo: "bar"}),
	)
	AssertNoError(t, h.Handle(context.Background(), rec))

	expected := fmt.Sprintf("%s INF foobar bool=true int=-12 uint=12 float=3.14 foo=bar time=%s dur=1s group.foo=bar group.subgroup.foo=bar err=the error stringer=stringer nostringer={bar}\r\n", now.Format(time.DateTime), now.Format(time.RFC3339))
	AssertEqual(t, expected, buf.String())
}

func TestHandler_WithAttr(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, &HandlerOptions{NoColor: true})
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", 0)
	h2 := h.WithAttrs([]slog.Attr{
		slog.Bool("bool", true),
		slog.Int("int", -12),
		slog.Uint64("uint", 12),
		slog.Float64("float", 3.14),
		slog.String("foo", "bar"),
		slog.Time("time", now),
		slog.Duration("dur", time.Second),
		slog.Group("group", slog.String("foo", "bar"), slog.Group("subgroup", slog.String("foo", "bar"))),
	})
	AssertNoError(t, h2.Handle(context.Background(), rec))

	expected := fmt.Sprintf("%s INF foobar bool=true int=-12 uint=12 float=3.14 foo=bar time=%s dur=1s group.foo=bar group.subgroup.foo=bar\r\n", now.Format(time.DateTime), now.Format(time.RFC3339))
	AssertEqual(t, expected, buf.String())

	buf.Reset()
	AssertNoError(t, h.Handle(context.Background(), rec))
	AssertEqual(t, fmt.Sprintf("%s INF foobar\r\n", now.Format(time.DateTime)), buf.String())
}

func TestHandler_WithGroup(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, &HandlerOptions{NoColor: true})
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", 0)
	rec.Add("int", 12)
	h2 := h.WithGroup("group1").WithAttrs([]slog.Attr{slog.String("foo", "bar")})
	AssertNoError(t, h2.Handle(context.Background(), rec))
	expected := fmt.Sprintf("%s INF foobar group1.foo=bar group1.int=12\r\n", now.Format(time.DateTime))
	AssertEqual(t, expected, buf.String())
	buf.Reset()

	h3 := h2.WithGroup("group2")
	AssertNoError(t, h3.Handle(context.Background(), rec))
	expected = fmt.Sprintf("%s INF foobar group1.foo=bar group1.group2.int=12\r\n", now.Format(time.DateTime))
	AssertEqual(t, expected, buf.String())

	buf.Reset()
	AssertNoError(t, h.Handle(context.Background(), rec))
	AssertEqual(t, fmt.Sprintf("%s INF foobar int=12\r\n", now.Format(time.DateTime)), buf.String())
}

func TestHandler_Levels(t *testing.T) {
	levels := map[slog.Level]string{
		slog.LevelDebug - 1: "DBG-1",
		slog.LevelDebug:     "DBG",
		slog.LevelDebug + 1: "DBG+1",
		slog.LevelInfo:      "INF",
		slog.LevelInfo + 1:  "INF+1",
		slog.LevelWarn:      "WRN",
		slog.LevelWarn + 1:  "WRN+1",
		slog.LevelError:     "ERR",
		slog.LevelError + 1: "ERR+1",
	}

	for l := range levels {
		t.Run(l.String(), func(t *testing.T) {
			buf := bytes.Buffer{}
			h := NewHandler(&buf, &HandlerOptions{Level: l, NoColor: true})
			for ll, s := range levels {
				AssertEqual(t, ll >= l, h.Enabled(context.Background(), ll))
				now := time.Now()
				rec := slog.NewRecord(now, ll, "foobar", 0)
				if ll >= l {
					AssertNoError(t, h.Handle(context.Background(), rec))
					AssertEqual(t, fmt.Sprintf("%s %s foobar\r\n", now.Format(time.DateTime), s), buf.String())
					buf.Reset()
				}
			}
		})
	}
}

func TestHandler_Source(t *testing.T) {
	buf := bytes.Buffer{}
	h := NewHandler(&buf, &HandlerOptions{NoColor: true, AddSource: true})
	h2 := NewHandler(&buf, &HandlerOptions{NoColor: true, AddSource: false})
	pc, file, line, _ := runtime.Caller(0)
	now := time.Now()
	rec := slog.NewRecord(now, slog.LevelInfo, "foobar", pc)
	AssertNoError(t, h.Handle(context.Background(), rec))
	cwd, _ := os.Getwd()
	file, _ = filepath.Rel(cwd, file)
	AssertEqual(t, fmt.Sprintf("%s INF %s:%d > foobar\r\n", now.Format(time.DateTime), file, line), buf.String())
	buf.Reset()
	AssertNoError(t, h2.Handle(context.Background(), rec))
	AssertEqual(t, fmt.Sprintf("%s INF foobar\r\n", now.Format(time.DateTime)), buf.String())
}

func TestHandler_Err(t *testing.T) {
	w := writerFunc(func(b []byte) (int, error) { return 0, errors.New("nope") })
	h := NewHandler(w, &HandlerOptions{NoColor: true})
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "foobar", 0)
	AssertError(t, h.Handle(context.Background(), rec))
}
