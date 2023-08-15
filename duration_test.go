package console

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	times := []time.Duration{
		2*time.Hour + 3*time.Minute + 4*time.Second + 5*time.Millisecond + 6*time.Microsecond + 7*time.Nanosecond,
		3*time.Minute + 4*time.Second + 5*time.Millisecond + 6*time.Microsecond + 7*time.Nanosecond,
		4*time.Second + 5*time.Millisecond + 6*time.Microsecond + 7*time.Nanosecond,
		5*time.Millisecond + 6*time.Microsecond + 7*time.Nanosecond,
		6*time.Microsecond + 7*time.Nanosecond,
		7 * time.Nanosecond,
		time.Duration(0),

		2*time.Hour + 7*time.Nanosecond,
		-2*time.Hour + 7*time.Nanosecond,
	}

	b := [4096]byte{}
	for _, tm := range times {
		bd := appendDuration(b[:0], tm)
		assert.Equal(t, tm.String(), string(bd))
	}

	bd := appendDuration(b[:0], 49*time.Hour+1*time.Second)
	assert.Equal(t, "2d1h0m1s", string(bd))
}

func BenchmarkDuration(b *testing.B) {
	d := 12*time.Hour + 13*time.Minute + 43*time.Second + 12*time.Millisecond
	b.Run("std", func(b *testing.B) {
		w := new(bytes.Buffer)
		w.Grow(2048)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w.WriteString(d.String())
			w.Reset()
		}
	})

	b.Run("append", func(b *testing.B) {
		w := new(buffer)
		w.Grow(2048)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w.AppendDuration(d)
			w.Reset()
		}
	})
}
