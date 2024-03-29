package system_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registration / Verification Tests", func() {
	It("responds to GET /registration with a status", func() {
		mobileNumber := "4169962646"
		verificationCode := "12345"
		pinCode := "12345"
		url := fmt.Sprintf(
			"http://%s/verify?mobile_number=%s&mobile_code=%s&mobile_pin=%s", serverAddress, mobileNumber, verificationCode, pinCode)
		resp, err := http.Get(url)

		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()

		// TODO : elaborate testing required
		Expect(resp.StatusCode).To(Equal(200))
	})
})
