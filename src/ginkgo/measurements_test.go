package ginkgo_test

import (
	"log"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Measurements", func() {
	log.SetOutput(GinkgoWriter)
	const N = 100

	Measure("Clone of testing benchmark", func(b Benchmarker) {
		var count int
		runtime := b.Time("runtime", func() {
			for i := 0; i < N; i++ {
				count++
				time.Sleep(10 * time.Millisecond)
			}
		})

		Expect(runtime.Seconds()).Should(BeNumerically("<", 3.0), "This thing shouldn't take too long.")

		b.RecordValue("Counts reached", float64(count))
	}, 5)

	Context("goroutines", func() {
		var queue chan struct{}
		var worker = func(count *int, queue chan struct{}, wg *sync.WaitGroup) {
			for _ = range queue {
				*count++
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			}
			wg.Done()
		}

		BeforeEach(func() {
			queue = make(chan struct{}, N)
			for i := 0; i < N; i++ {
				queue <- struct{}{}
			}
			close(queue)
		})

		Measure("Message throuput with one goroutine", func(b Benchmarker) {
			var countA int

			runtime := b.Time("runtime", func() {
				wg := sync.WaitGroup{}
				go worker(&countA, queue, &wg)
				wg.Wait()
			})

			Expect(runtime.Seconds()).Should(BeNumerically("<", 3.0), "This thing shouldn't take too long.")

			b.RecordValue("countA", float64(countA))
		}, 5)

		Measure("Message throuput with two goroutines", func(b Benchmarker) {
			var countA int
			var countB int

			runtime := b.Time("runtime", func() {
				wg := sync.WaitGroup{}
				wg.Add(2)
				go worker(&countA, queue, &wg)
				go worker(&countB, queue, &wg)
				wg.Wait()
			})

			Expect(runtime.Seconds()).Should(BeNumerically("<", 3.0), "This thing shouldn't take too long.")

			b.RecordValue("countA", float64(countA))
			b.RecordValue("countB", float64(countB))
		}, 5)
	})
})
