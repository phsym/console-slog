package console

import (
	"bytes"
	"testing"
)

func AssertZero[E comparable](t *testing.T, v E) {
	t.Helper()
	var zero E
	if v != zero {
		t.Errorf("expected zero value, got %v", v)
	}
}

func AssertEqual[E comparable](t *testing.T, expected, value E) {
	t.Helper()
	if expected != value {
		t.Errorf("expected %v, got %v", expected, value)
	}
}

func TestBuffer(t *testing.T) {
	b := new(buffer)
	AssertZero(t, b.Len())
	b.AppendString("foobar")
	AssertEqual(t, 6, b.Len())
	b.AppendString("baz")
	AssertEqual(t, 9, b.Len())
	AssertEqual(t, "foobarbaz", b.String())

	b.AppendByte('.')
	AssertEqual(t, 10, b.Len())
	AssertEqual(t, "foobarbaz.", b.String())
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
