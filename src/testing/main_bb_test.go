package testing_test

// If we define package this way, it is considered separate from our main package
// and only has access to exported (i.e. capitalized) members.

// The file name needs to end with "_test.go"

// Run this with `go test`
// Run `go test -v` for verbose output
// Also `go test -help` is your friend
// https://golang.org/pkg/testing/

import (
	"testing"
	"time"

	. "github.com/egurnov/ginkgo-talk/src/testing"
)

// Each exported function starting with "Test" is considered a test case
func TestFails(t *testing.T) {
	t.Fail()
}

func TestSkip(t *testing.T) {
	t.Skip("Who needs tests anyway?")
}

func TestReturnsTrue(t *testing.T) {
	res := ReturnTrue()
	if res != true {
		t.Errorf("Error: returned %v instead of %v", res, true)
	}
}

func TestHierarchy(t *testing.T) {
	t.Run("subtest1", func(t *testing.T) {
		t.Run("subsubtest1", func(t *testing.T) {})
	})
	t.Run("subtest2", func(t *testing.T) {})
	t.Run("subtest3", func(t *testing.T) {})
}

// Each exported function starting with "Benchmark" is considered a benchmark.
// Benchmarks are not run automatically. Use `go test -test.bench=.` to run them all
// or provide a regexp to run specific benchmarks.
func BenchmarkAddition(b *testing.B) {
	time.Sleep(1 * time.Second) // Do some expensive setup
	count := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count++
		time.Sleep(10 * time.Millisecond)
	}
}
