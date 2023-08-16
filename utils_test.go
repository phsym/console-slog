package console

import (
	"cmp"
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

func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

// func AssertNil(t *testing.T, value any) {
// 	t.Helper()
// 	if value != nil {
// 		t.Errorf("expected nil, got %v", value)
// 	}
// }

type writerFunc func([]byte) (int, error)

func (w writerFunc) Write(b []byte) (int, error) {
	return w(b)
}
