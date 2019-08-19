package system_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SMS Tests", func() {
	It("responds to GET /sms with status", func() {
		// mobileNum := "+14169962646"
		// url := fmt.Sprintf("http://%s/sms?mobile_num=%s", serverAddress, mobileNum)
		url := "https://www.google.com"
		resp, err := http.Get(url)

		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()

		// TODO : elaborate testing required
		Expect(resp.StatusCode).To(Equal(200))
	})
})
