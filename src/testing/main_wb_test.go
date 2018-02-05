package testing

import "testing"

// Tests in this file are part of main package, so they have access to non-exported members.
func TestReturnsFalse(t *testing.T) {
	res := returnFalse()
	if res != false {
		t.Errorf("Error: returned %v instead of %v", res, true)
	}
}
