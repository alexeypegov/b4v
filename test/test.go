package test

import (
  "testing"
  "strconv"
)

// AssertEquals compares two strings
func AssertEquals(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nE: %s\nA: %s", strconv.Quote(expected), strconv.Quote(actual))
	}
}