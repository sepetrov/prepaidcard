// Package testing has test helper functions.
package testing

import "testing"

// Must is a test helper, which interrupts test t and printfs s with args if ok == false.
func Must(t *testing.T, ok bool, s string, args ...interface{}) {
	t.Helper()
	if !ok {
		t.Fatalf(s, args...)
	}
}

// MustE os a test helper, which interrupts test t if a != b.
func MustE(t *testing.T, a interface{}, b interface{}, s string, args ...interface{}) {
	if len(s) == 0 {
		s = "got %#v != %#v, want a == b"
	}
	p := []interface{}{a, b}
	if len(args) > 0 {
		p = append(p, args...)
	}
	Must(t, a == b, s, p...)
}

// MustErr is a test helper, which interrupts test t and printfs s with args if err == nil.
func MustErr(t *testing.T, err error, s string, args ...interface{}) {
	t.Helper()
	if err == nil {
		t.Fatalf(s, args...)
	}
}

// MustNotErr is a test helper, which interrupts test t and printfs s with args if err != nil.
func MustNotErr(t *testing.T, err error, s string, args ...interface{}) {
	t.Helper()
	if err != nil {
		args := append(args, err)
		t.Fatalf(s, args...)
	}
}
