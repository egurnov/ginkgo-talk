package ginkgo_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Httptest", func() {
	Context("httptest", func() {
		var server *httptest.Server
		var response *http.Response
		var err error

		BeforeEach(func() {
			server = httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Hello, World!"))
				}),
			)
		})

		JustBeforeEach(func() {
			response, err = http.Get(server.URL)
		})

		AfterEach(func() {
			server.Close()
		})

		It("works", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body).To(ContainSubstring("Hello"))
		})
	})
})
