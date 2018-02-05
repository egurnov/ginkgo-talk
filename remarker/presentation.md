name: Testing in Go
class: center, middle

# Testing in Go

### Introduction to testing package, Ginkgo and Gomega

Alexander Egurnov, IBM  
5 February 2018

---

# Overview

* testing package
* Ginkgo & Gomega
* httptest
* code generation

---

class: center, middle
# testing

---

# testing package

* built into language
* run with `go test`
* test file name should end with "`_test.go`"
* runs all functions starting with "`Test`"

```go
func TestMe(t *testing.T) {
	// ...
}
```

* no assertions, need to use `t.Fail(), t.Error(), t.Fatal()` and friends
* can skip tests with `t.Skip()`
* package typically gets a `_test` suffix

---
# Benchmarks

* Run when `-bench` flag is given.
* Should repeat the action b.N times.
* Will run enough times to be running for 1 second.

```go
func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
```

???

`go test -v -run=X  -test.bench=.`

---

#Subtests

Allow to share setup and teardown, create test hierarchies.

```go
func TestHierarchy(t *testing.T) {
	// <setup code>
	t.Run("subtest1", func(t *testing.T) {
		t.Run("subsubtest1", func(t *testing.T){ ... })
	})
	t.Run("subtest2", func(t *testing.T) { ... })
	t.Run("subtest3", func(t *testing.T) { ... })
	// <tear-down code>
}
```

???

`go test -v -run=Hierarchy`

---

class: center, middle
# Ginkgo & Gomega

---

# Ginkgo

BDD-style testing framework

Can be paired with any assertions library, although Gomega is preferred.

---

# Getting started

```console
$ ginkgo bootstrap
```
creates a test suite file, which plugs into standard testing framework.

```go
package ginkgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail) // The only connection between Ginkgo and Gomega
	RunSpecs(t, "Ginkgo Suite")
}
```

Generate some boilerplate to start writing actual tests.
```console
$ ginkgo generate Subject
```

---

# Ginkgo DSL

* `Describe` and `Context` to introduce test groups.
* `It` and `Specify` to introduce a single spec.
* `BeforeEach`, `AfterEach`, `BeforeSuite`, `AfterSuite` for shared setup and teardown.
* Use closures to share state between functions.
* `Fail("Failure reason")` to fail a test.
* Can write to `GinkgoWriter` to hide test output, unless tests are stopped of failed.
```go
log.SetOutput(GinkgoWriter)
```

---

# Running tests

* Run tests with
```console
$ ginkgo
```

* Additionally traverse all subdirectories with
```console
$ ginkgo -r
```

* Run ginkgo on every save.
```console
$ ginkgo watch -race -notify
```
	* Can send notifications with terminal-notifier.  
		Need to `brew install terminal-notifier`

---

# Gomega DSL

* Two notations:
```go
Expect(ACTUAL).To(Equal(EXPECTED))
Ω(ACTUAL).Should(Equal(EXPECTED))
```

* Implicitly checks for errors
	* all extra returns should be nil or zero

* Annotations
```go
Expect(ACTUAL).To(Equal(EXPECTED), "My annotation %d", foo)
```

---

# Matchers

```go
Equal(interface{}) // deep equality
BeEquivalentTo(interface{}) // casts to type of EXPECTED
BeIdenticalTo(interface{}) // can compare pointers

BeNumerically(comp string, interface{}, [threshold time.Duration])
BeTemporarily(comp string, interface{}, [threshold time.Duration])

BeNil()
BeZero() // zero value for its type or nil.
BeTrue()
BeFalse()

HaveOccurred()
Succeeded()
MatchError(interface{})

BeClosed()
Receive([&something])
Receive(Equal("this"))
BeSent(interface{})
```

---

# Matchers (cont.)

```go
BeAnExistingFile() // ACTUAL should be a file path
BeARegularFile()
BeADirectory()

ContainSubstring(string)
HavePrefix(string)
HaveSuffix(string)
MatchRegexp(regexp)
MatchJSON(interface{}) // string, []byte or Stringer
MatchUnorderedJSON(interface{})
MatchXML(interface{})
MatchYAML(interface{})

BeEmpty()
HaveLen(int)
HaveCap(int)
ContainElement(interface{})
ConsistOf(interface{})
HaveKey(interface{})
HaveKeyWithValue(interface{})

Panic()
```

---

# Matchers (cont.)

```go
And(...)
Or(...)
SatisfyAll(...)
SatisfyAny(...)
Not(matcher)

WithTransform(transform func(), matcher)
```

And you can write your own matchers, if this is not enough.

---

# Skipping and Focusing

* Prevent specs or spec groups from running by prepending `X` or `P` it, i.e. `XIt` or `PContext`.
 * Or call `Skip("reason")` at runtime.
* Focus by prepending "F".
* Can use `--focus=REGEXP` and `--skip=REGEXP` from command line.
* `ginkgo unfocus` to remove all `F`s from source files.

---

# Complex setup scenarios

* Use `JustBeforeEach` to separate creation from configuration.

```go
Describe("Outer", func() {
	var param string

	BeforeEach(func() {
		param = "this"
	})
	JustBeforeEach(func() {
		DoStuff(param)
	})

	Context("Inner", func() {
		BeforeEach(func() {
			param = "that"
		})
	})
})
```

---

# Measurements

* Can run measurements with `Measure`.
* Only run the given number of times.
* Can record additional values.
* Basic statistics will be reported. Can provide custom reporters.
* Can be skipped: `ginkgo --skipMeasurements`

```go
Measure("Name", func(b Benchmarker) {
	runtime := b.Time("runtime", func() {
		DoStuff()
	})

	Expect(runtime.Seconds()).Should(BeNumerically("<", 3.0), "This thing shouldn't take too long.")

	b.RecordValue("Counts reached", float64(count))
}, 5)
```

---

# Asynchronous tests

* All non-container function can receive an optional signal channel argument.
* Optionally can set a timeout in seconds

```go
It("can handle async functions", func(done Done) {
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()
}, 10)
```

* Use `defer GinkgoRecover()` when asserting in goroutines.

---

# Eventually / Consistently

* `Eventually` polls the argument until it succeeds. `Consistently` checks that assertion holds (at least) for a period of time.
* Timeout and polling interval are optional and can be `time.Duration`, `float64` (seconds) or strings, e.g. "200ms".

```go
Eventually(func() []int {
	return thing.SliceImMonitoring
}, TIMEOUT, POLLING_INTERVAL).Should(HaveLen(2))

Consistently(func() []int {
	return thing.MemoryUsage()
}, DURATION, POLLING_INTERVAL).Should(BeNumerically("<", 10))
```

* Work well with channels.

```go
Eventually(channel).Should(BeClosed())
Eventually(channel).Should(Receive())
```

---

class: center, middle
# httptest

---

# httptest

```go
server = httptest.NewServer(
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}),
)
```

```go
response, err = http.Get(server.URL)
```

---

class: center, middle
# Stubs, mocks, fakes

---

# Stubs, mocks, fakes

* Write everything yourself
* `//go:generate counterfeiter . TargetInterface`
* ...
