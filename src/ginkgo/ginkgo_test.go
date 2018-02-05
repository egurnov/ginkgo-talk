package ginkgo_test

import (
	"errors"
	"log"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ginkgo", func() {

	XIt("always fails", func() {
		Fail("Noooooooo!")
	})

	Describe("BeforeEach/JustBeforeEach/AfterEach", func() {
		BeforeEach(func() {
			log.Print("BeforeEach 1")
		})

		JustBeforeEach(func() {
			log.Print("JustBeforeEach 1")
		})

		AfterEach(func() {
			log.Print("AfterEach 1")
		})

		Context("inner context", func() {
			BeforeEach(func() {
				log.Print("BeforeEach 2")
			})

			JustBeforeEach(func() {
				log.Print("JustBeforeEach 2")
			})

			AfterEach(func() {
				log.Print("AfterEach 2")
			})

			It("runs the test", func() {
				log.Print("It")
			})
		})
	})

	Context("implicitly checks for errors", func() {
		It("allows nil/zero values", func() {
			good := func() (string, error) {
				return "OK", nil
			}
			Expect(good()).To(Equal("OK"))
		})

		XIt("fails on non-nil/non-zero extra returns", func() {
			bad := func() (string, error) {
				return "OK", errors.New("Ooops!")
			}
			Expect(bad()).To(Equal("OK"))
		})
	})

	It("works and showcases Gomega matchers", func() {
		Expect(5).To(Equal(5))    // Deep equality, strict about types
		Expect(5).NotTo(Equal(3)) // Also with negations
		Expect(5).ToNot(Equal(3)) //
		Ω(5).Should(Equal(5))     // This is just another syntax
		Expect(5).To(Equal(5),
			"Basic math should work as expected") // Also with anotations

		type KindOfInt int
		const (
			Zero KindOfInt = iota
			One
			Two
		)

		Expect(Zero).To(BeEquivalentTo(0)) // Performs type casting
		Expect(Zero).ToNot(Equal(0))
		Expect(5.1).To(BeEquivalentTo(5)) // Type casting gotcha
		Expect(5).ToNot(BeEquivalentTo(5.1))
		Expect(map[string]int{"a": 1, "b": 2}).
			To(BeEquivalentTo(map[string]int{"b": 2, "a": 1}))

		p1 := &struct{ v int }{v: 5}
		p2 := p1
		Expect(p1).To(BeIdenticalTo(p2))
		Expect(p1).ToNot(BeIdenticalTo(&struct{ v int }{v: 5}))
		Expect(p1).ToNot(BeIdenticalTo(struct{ v int }{v: 5}))

		Expect(5).To(BeNumerically("<", 5.1))         // The rigth way
		Expect(5).To(BeNumerically("~", 5.005, 1e-2)) // Comparison with a threshold
		d1 := time.Date(2018, time.February, 5, 19, 30, 0, 0, time.UTC)
		d2 := time.Date(2018, time.February, 5, 19, 34, 0, 0, time.UTC)
		Expect(d1).To(BeTemporally("~", d2, 5*time.Minute))

		var p *int
		Expect(p).To(BeNil())
		Expect("").To(BeZero())

		_, err := strconv.Atoi("not a number")
		Expect(err).To(HaveOccurred())
		_, err = strconv.Atoi("42")
		Expect(err).To(Succeed())

		ch := make(chan int, 1)
		var v int
		ch <- 5
		Expect(ch).To(Receive(&v))
		Expect(ch).To(BeSent(7))
		Expect(ch).To(Receive(Equal(7)))

		Expect("Golang").To(HavePrefix("Go"))
		Expect("Abracadabra").To(ContainSubstring("cad"))
		Expect("x-y=z").To(ContainSubstring("%v-%v", "x", "y"))
		Expect("{\"a\": 1, \"b\": 2}").To(MatchJSON("{\"b\": 2, \"a\": 1}"))

		// Collections
		theSequence := []int{4, 8, 15, 16, 23, 42}
		Expect(theSequence).ToNot(BeEmpty())
		Expect(theSequence).To(HaveLen(6))
		Expect(theSequence).To(ContainElement(23))
		Expect(theSequence).To(ConsistOf(8, 16, 42, 23, 15, 4))

		shoppingList := map[string]int{"apples": 4, "tomatoes": 10, "milk": 1}
		Expect(shoppingList).To(HaveKey("apples"))
		Expect(shoppingList).To(HaveKeyWithValue("tomatoes", 10))

		Expect(func() { panic("Bummer") }).To(Panic())

		Expect(5).To(
			And( // Optionally SatisfyAll()
				BeNumerically(">", 4),
				BeNumerically("<", 6),
			),
		)
		Expect(5).To(
			Or( // Optionally SatisfyAny()
				BeNumerically(">", 0),
				BeNumerically("<", 0),
			),
		)
		Expect(5).To(Not(BeNil())) // Can also negate a single matcher
	})

	Describe("Asynchronous functions", func() {
		It("can handle async functions", func(done Done) {
			go func() {
				time.Sleep(3 * time.Second)
				close(done)
			}()
		}, 5)

		XIt("can handle async function failures", func(done Done) {
			go func() {
				defer GinkgoRecover()

				time.Sleep(1 * time.Second)
				Fail("You shall not pass!")
				close(done)
			}()
		}, 3)
	})
})
