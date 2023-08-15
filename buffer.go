package console

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"time"
)

type buffer []byte

func (b *buffer) Grow(n int) {
	*b = slices.Grow(*b, n)
}

func (b *buffer) Bytes() []byte {
	return *b
}

func (b *buffer) String() string {
	return string(*b)
}

func (b *buffer) Len() int {
	return len(*b)
}

func (b *buffer) WriteTo(dst io.Writer) (int64, error) {
	l := len(*b)
	if l == 0 {
		return 0, nil
	}
	n, err := dst.Write(*b)
	if err != nil {
		return int64(n), err
	}
	if n < l {
		return int64(n), io.ErrShortWrite
	}
	b.Reset()
	return int64(n), nil
}

func (b *buffer) Reset() {
	*b = (*b)[:0]
}

func (b *buffer) Clone() buffer {
	return append(buffer(nil), *b...)
}

func (b *buffer) Clip() {
	*b = slices.Clip(*b)
}

func (buf *buffer) copy(src *buffer) {
	if src.Len() > 0 {
		buf.Append(src.Bytes())
	}
}

func (b *buffer) Append(data []byte) {
	*b = append(*b, data...)
}

func (b *buffer) AppendString(s string) {
	*b = append(*b, s...)
}

// func (b *buffer) AppendQuotedString(s string) {
// 	b.buff = strconv.AppendQuote(b.buff, s)
// }

func (b *buffer) AppendByte(byt byte) {
	*b = append(*b, byt)
}

func (b *buffer) AppendTime(t time.Time, format string) {
	*b = t.AppendFormat(*b, format)
}

func (b *buffer) AppendInt(i int64) {
	*b = strconv.AppendInt(*b, i, 10)
}

func (b *buffer) AppendUint(i uint64) {
	*b = strconv.AppendUint(*b, i, 10)
}

func (b *buffer) AppendFloat(i float64) {
	*b = strconv.AppendFloat(*b, i, 'g', -1, 64)
}

func (b *buffer) AppendBool(i bool) {
	*b = strconv.AppendBool(*b, i)
}

func (b *buffer) AppendDuration(d time.Duration) {
	*b = appendDuration(*b, d)
}

func (buf *buffer) NewLine() {
	buf.AppendString("\r\n")
}

func (b *buffer) withColor(c color, f func()) {
	if c == "" {
		f()
		return
	}
	b.AppendString(string(c))
	f()
	b.AppendString(string(reset))
}

func (w *buffer) writeColoredTime(t time.Time, format string, seq color) {
	w.withColor(seq, func() {
		w.AppendTime(t, format)
	})
}

func (w *buffer) writeColoredString(s string, seq color) {
	w.withColor(seq, func() {
		w.AppendString(s)
	})
}

func (w *buffer) writeColoredInt(i int64, seq color) {
	w.withColor(seq, func() {
		w.AppendInt(i)
	})
}

func (w *buffer) writeColoredUint(i uint64, seq color) {
	w.withColor(seq, func() {
		w.AppendUint(i)
	})
}

func (w *buffer) writeColoredFloat(i float64, seq color) {
	w.withColor(seq, func() {
		w.AppendFloat(i)
	})
}

func (w *buffer) writeColoredBool(b bool, seq color) {
	w.withColor(seq, func() {
		w.AppendBool(b)
	})
}

func (w *buffer) writeColoredDuration(d time.Duration, seq color) {
	w.withColor(seq, func() {
		w.AppendDuration(d)
	})
}

func (buf *buffer) writeTimestamp(tt time.Time) {
	buf.writeColoredTime(tt, time.DateTime, colorTimestamp)
	buf.AppendByte(' ')
}

func (buf *buffer) writeSource(pc uintptr, cwd string) {
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	if cwd != "" {
		if ff, err := filepath.Rel(cwd, frame.File); err == nil {
			frame.File = ff
		}
	}
	buf.withColor(colorSource, func() {
		buf.AppendString(frame.File)
		buf.AppendByte(':')
		buf.AppendInt(int64(frame.Line))
	})
	buf.writeColoredString(" > ", colorAttrKey)
}

func (buf *buffer) writeMessage(msg string) {
	buf.writeColoredString(msg, colorMessage)
}

func (buf *buffer) writeAttr(a slog.Attr, group string) {
	value := a.Value.Resolve()
	if value.Kind() == slog.KindGroup {
		subgroup := a.Key
		if group != "" {
			subgroup = group + "." + a.Key
		}
		for _, attr := range value.Group() {
			buf.writeAttr(attr, subgroup)
		}
		return
	}
	buf.AppendByte(' ')
	buf.withColor(colorAttrKey, func() {
		if group != "" {
			buf.AppendString(group)
			buf.AppendByte('.')
		}
		buf.AppendString(a.Key)
		buf.AppendByte('=')
	})
	buf.writeValue(value)
	return
}

func (buf *buffer) writeValue(value slog.Value) {
	switch value.Kind() {
	case slog.KindInt64:
		buf.writeColoredInt(value.Int64(), colorAttrValue)
	case slog.KindBool:
		buf.writeColoredBool(value.Bool(), colorAttrValue)
	case slog.KindFloat64:
		buf.writeColoredFloat(value.Float64(), colorAttrValue)
	case slog.KindTime:
		buf.writeColoredTime(value.Time(), time.RFC3339, colorAttrValue)
	case slog.KindUint64:
		buf.writeColoredUint(value.Uint64(), colorAttrValue)
	case slog.KindDuration:
		buf.writeColoredDuration(value.Duration(), colorAttrValue)
	case slog.KindAny:
		switch v := value.Any().(type) {
		case error:
			buf.writeColoredString(v.Error(), colorErrorValue)
			return
		case fmt.Stringer:
			buf.writeColoredString(v.String(), colorAttrValue)
		}
		fallthrough
	// case slog.KindString:
	// 	fallthrough
	default:
		buf.writeColoredString(value.String(), colorAttrValue)
	}
}

func (buf *buffer) writeLevel(l slog.Level) {
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
	buf.writeColoredString(str, style)
	buf.AppendByte(' ')
}
