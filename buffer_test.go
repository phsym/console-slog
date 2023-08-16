package console

import (
	"bytes"
	"cmp"
	"testing"
	"time"
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

func AssertNotEqual[E comparable](t *testing.T, expected, value E) {
	t.Helper()
	if expected == value {
		t.Errorf("expected to be different, got %v", value)
	}
}

func AssertGreaterOrEqual[E cmp.Ordered](t *testing.T, expected, value E) {
	t.Helper()
	if expected > value {
		t.Errorf("expected to be %v to be greater than %v", value, expected)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
	}
}

// func AssertNil(t *testing.T, value any) {
// 	t.Helper()
// 	if value != nil {
// 		t.Errorf("expected nil, got %v", value)
// 	}
// }

func TestBuffer_Append(t *testing.T) {
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

	b.AppendBool(true)
	b.AppendBool(false)
	b.AppendFloat(3.14)
	b.AppendInt(42)
	b.AppendUint(12)
	b.Append([]byte("foo"))
	b.AppendDuration(1 * time.Second)
	now := time.Now()
	b.AppendTime(now, time.RFC3339)

	AssertEqual(t, "foobarbaz.truefalse3.144212foo1s"+now.Format(time.RFC3339), b.String())
}

func TestBuffer_WriteTo(t *testing.T) {
	dest := bytes.Buffer{}
	b := new(buffer)
	n, err := b.WriteTo(&dest)
	AssertNoError(t, err)
	AssertZero(t, n)
	b.AppendString("foobar")
	n, err = b.WriteTo(&dest)
	AssertEqual(t, len("foobar"), int(n))
	AssertNoError(t, err)
	AssertEqual(t, "foobar", dest.String())
	AssertZero(t, b.Len())
}

func TestBuffer_Clone(t *testing.T) {
	b := new(buffer)
	b.AppendString("foobar")
	b2 := b.Clone()
	AssertEqual(t, b.String(), b2.String())
	AssertNotEqual(t, &b.Bytes()[0], &b2.Bytes()[0])
}

func TestBuffer_Copy(t *testing.T) {
	b := new(buffer)
	b.AppendString("foobar")
	b2 := new(buffer)
	b2.copy(b)
	AssertEqual(t, b.String(), b2.String())
	AssertNotEqual(t, &b.Bytes()[0], &b2.Bytes()[0])
}

func TestBuffer_Reset(t *testing.T) {
	b := new(buffer)
	b.AppendString("foobar")
	AssertEqual(t, "foobar", b.String())
	AssertEqual(t, len("foobar"), b.Len())
	bufCap := b.Cap()
	b.Reset()
	AssertZero(t, b.Len())
	AssertEqual(t, bufCap, b.Cap())
}

func TestBuffer_Grow(t *testing.T) {
	b := new(buffer)
	AssertZero(t, b.Cap())
	b.Grow(12)
	AssertGreaterOrEqual(t, 12, b.Cap())
	b.Grow(6)
	AssertGreaterOrEqual(t, 12, b.Cap())
	b.Grow(24)
	AssertGreaterOrEqual(t, 24, b.Cap())
}

func TestBuffer_Clip(t *testing.T) {
	b := new(buffer)
	b.AppendString("foobar")
	b.Grow(12)
	AssertGreaterOrEqual(t, 12, b.Cap())
	b.Clip()
	AssertEqual(t, "foobar", b.String())
	AssertEqual(t, len("foobar"), b.Cap())
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
