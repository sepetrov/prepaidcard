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
