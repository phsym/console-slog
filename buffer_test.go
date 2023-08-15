package console

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	b := new(buffer)
	assert.Zero(t, b.Len())
	b.AppendString("foobar")
	assert.Equal(t, 6, b.Len())
	b.AppendString("baz")
	assert.Equal(t, 9, b.Len())
	assert.Equal(t, "foobarbaz", b.String())

	b.AppendByte('.')
	assert.Equal(t, 10, b.Len())
	assert.Equal(t, "foobarbaz.", b.String())
}

func BenchmarkBuffer(b *testing.B) {
	data := []byte("foobarbaz")

	b.Run("std", func(b *testing.B) {
		buf := bytes.Buffer{}
		for i := 0; i < b.N; i++ {
			buf.Write(data)
			buf.WriteByte('.')
			buf.Reset()
		}
	})

	b.Run("buffer", func(b *testing.B) {
		buf := buffer{}
		for i := 0; i < b.N; i++ {
			buf.Append(data)
			buf.AppendByte('.')
			buf.Reset()
		}
	})
}
