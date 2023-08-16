package console

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

type encoder struct {
	nocolor bool
}

func (e *encoder) NewLine(buf *buffer) {
	buf.AppendString("\r\n")
}

func (e *encoder) withColor(b *buffer, c color, f func()) {
	if c == "" || e.nocolor {
		f()
		return
	}
	b.AppendString(string(c))
	f()
	b.AppendString(string(reset))
}

func (e *encoder) writeColoredTime(w *buffer, t time.Time, format string, seq color) {
	e.withColor(w, seq, func() {
		w.AppendTime(t, format)
	})
}

func (e *encoder) writeColoredString(w *buffer, s string, seq color) {
	e.withColor(w, seq, func() {
		w.AppendString(s)
	})
}

func (e *encoder) writeColoredInt(w *buffer, i int64, seq color) {
	e.withColor(w, seq, func() {
		w.AppendInt(i)
	})
}

func (e *encoder) writeColoredUint(w *buffer, i uint64, seq color) {
	e.withColor(w, seq, func() {
		w.AppendUint(i)
	})
}

func (e *encoder) writeColoredFloat(w *buffer, i float64, seq color) {
	e.withColor(w, seq, func() {
		w.AppendFloat(i)
	})
}

func (e *encoder) writeColoredBool(w *buffer, b bool, seq color) {
	e.withColor(w, seq, func() {
		w.AppendBool(b)
	})
}

func (e *encoder) writeColoredDuration(w *buffer, d time.Duration, seq color) {
	e.withColor(w, seq, func() {
		w.AppendDuration(d)
	})
}

func (e *encoder) writeTimestamp(buf *buffer, tt time.Time) {
	e.writeColoredTime(buf, tt, time.DateTime, colorTimestamp)
	buf.AppendByte(' ')
}

func (e *encoder) writeSource(buf *buffer, pc uintptr, cwd string) {
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	if cwd != "" {
		if ff, err := filepath.Rel(cwd, frame.File); err == nil {
			frame.File = ff
		}
	}
	e.withColor(buf, colorSource, func() {
		buf.AppendString(frame.File)
		buf.AppendByte(':')
		buf.AppendInt(int64(frame.Line))
	})
	e.writeColoredString(buf, " > ", colorAttrKey)
}

func (e *encoder) writeMessage(buf *buffer, msg string) {
	e.writeColoredString(buf, msg, colorMessage)
}

func (e *encoder) writeAttr(buf *buffer, a slog.Attr, group string) {
	value := a.Value.Resolve()
	if value.Kind() == slog.KindGroup {
		subgroup := a.Key
		if group != "" {
			subgroup = group + "." + a.Key
		}
		for _, attr := range value.Group() {
			e.writeAttr(buf, attr, subgroup)
		}
		return
	}
	buf.AppendByte(' ')
	e.withColor(buf, colorAttrKey, func() {
		if group != "" {
			buf.AppendString(group)
			buf.AppendByte('.')
		}
		buf.AppendString(a.Key)
		buf.AppendByte('=')
	})
	e.writeValue(buf, value)
	return
}

func (e *encoder) writeValue(buf *buffer, value slog.Value) {
	switch value.Kind() {
	case slog.KindInt64:
		e.writeColoredInt(buf, value.Int64(), colorAttrValue)
	case slog.KindBool:
		e.writeColoredBool(buf, value.Bool(), colorAttrValue)
	case slog.KindFloat64:
		e.writeColoredFloat(buf, value.Float64(), colorAttrValue)
	case slog.KindTime:
		e.writeColoredTime(buf, value.Time(), time.RFC3339, colorAttrValue)
	case slog.KindUint64:
		e.writeColoredUint(buf, value.Uint64(), colorAttrValue)
	case slog.KindDuration:
		e.writeColoredDuration(buf, value.Duration(), colorAttrValue)
	case slog.KindAny:
		switch v := value.Any().(type) {
		case error:
			e.writeColoredString(buf, v.Error(), colorErrorValue)
			return
		case fmt.Stringer:
			e.writeColoredString(buf, v.String(), colorAttrValue)
			return
		}
		fallthrough
	case slog.KindString:
		fallthrough
	default:
		e.writeColoredString(buf, value.String(), colorAttrValue)
	}
}

func (e *encoder) writeLevel(buf *buffer, l slog.Level) {
	var style color
	var str string
	switch {
	case l >= slog.LevelError:
		style = colorLevelError
		str = "ERR"
	case l >= slog.LevelWarn:
		style = colorLevelWarn
		str = "WRN"
	case l >= slog.LevelInfo:
		style = colorLevelInfo
		str = "INF"
	case l >= slog.LevelDebug:
		style = colorLevelDebug
		str = "DBG"
	default:
		style = bold
		str = "???"
	}
	e.writeColoredString(buf, str, style)
	buf.AppendByte(' ')
}
