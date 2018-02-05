package ginkgo_test

import (
	"errors"
	"fmt"

	. "github.com/egurnov/ginkgo-talk/src/handler"
	. "github.com/egurnov/ginkgo-talk/src/handler/handlerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var shortTimeout = 4.0

var _ = Describe("Handler", func() {

	var (
		fc  *FakeClient
		msg string
	)

	BeforeEach(func() {
		fc = &FakeClient{}
		msg = ""
	})

	Describe("Handler", func() {
		var (
			err error
			h   *RequestHandler
		)

		BeforeEach(func() {
			h = NewHandler("Scotty", fc)
		})

		JustBeforeEach(func() {
			err = h.Handle(msg)
		}, shortTimeout)

		Context("when msg is empty", func() {
			BeforeEach(func() {
				msg = ""
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when msg is not empty", func() {
			BeforeEach(func() {
				msg = "Hello"
			})

			It("works", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(fc.SendArgsForCall(0)).To(And(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})
		})

		Context("when asked a question", func() {
			BeforeEach(func() {
				msg = "Hello?"
			})

			It("is slow", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(fc.SendArgsForCall(0)).To(And(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})
		})
	})

	Describe("AsyncHandler", func() {
		var errChan chan error
		var h *AsyncHandler

		BeforeEach(func() {
			errChan = make(chan error)
			h = NewAsyncHandler(NewHandler("Scotty", fc), errChan)
		})

		JustBeforeEach(func() {
			h.Handle(msg)
		})

		Context("when msg is empty", func() {
			BeforeEach(func() {
				msg = ""
			})

			It("sends an error", func() {
				var err error
				Eventually(errChan).Should(Receive(&err))
				Expect(err).ToNot(BeNil())
			})
		})

		Context("when msg is not empty", func() {
			BeforeEach(func() {
				msg = "Hello"
			})

			It("works", func() {
				Eventually(fc.SendCallCount).Should(BeNumerically(">", 0))
				Expect(fc.SendArgsForCall(0)).To(And(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})

			It("works like this as well", func() {
				Eventually(func() (string, error) {
					if fc.SendCallCount() > 0 {
						return fc.SendArgsForCall(0), nil
					}
					return "", errors.New("Not yet")
				}).Should(SatisfyAll(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})

			It("works with assertions inside a goroutine", func() {
				syncChan := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					err := <-errChan
					fmt.Println("received result")
					Expect(err).To(BeNil())
					// Expect(err).ToNot(BeNil()) // will panic without GinkgoRecover
					close(syncChan)
					fmt.Println("closed sync")
				}()
				fmt.Println("waiting sync")
				<-syncChan
				fmt.Println("received sync")
				Expect(fc.SendArgsForCall(0)).To(And(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})
		})

		Context("when asked a question", func() {
			BeforeEach(func() {
				msg = "Hello?"
			})

			It("is still slow", func() {
				Eventually(errChan, shortTimeout).Should(Receive(BeNil()))
				Expect(fc.SendArgsForCall(0)).To(And(HaveSuffix("Scotty"), ContainSubstring(msg)))
			})
		})
	})
})
