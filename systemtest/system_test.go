package system_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the webserver", func() {
	It("responds to GET /auth with an http redirect to google signin url", func() {
		url := fmt.Sprintf("http://%s/auth", serverAddress)
		resp, err := http.Get(url)

		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()

		// TODO : elaborate testing required
		Expect(resp.StatusCode).To(Equal(200))
	})
	It("responds to GET /auth/callback with a google signin url", func() {
		url := fmt.Sprintf("http://%s/auth/callback", serverAddress)
		resp, err := http.Get(url)

		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()

		// TODO : elaborate testing required
		Expect(resp.StatusCode).To(Equal(403))
	})
})
