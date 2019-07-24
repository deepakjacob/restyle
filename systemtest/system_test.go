package system_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/deepakjacob/restyle/logger"
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
	It("/upload should respond with 200", func() {
		if err := logger.Init(-1, "2006-01-02T15:04:05.000Z"); err != nil {
			log.Fatal(err)
		}
		path, _ := os.Getwd()
		path += "/test_resources/img_1.jpg"
		formParams := map[string]string{
			"obj_type":       "saree",
			"material":       "silk",
			"speciality":     "kancheepuram",
			"dress_category": "Women",
			"age_min":        "18",
			"age_max":        "100",
			"tags":           `2019,onam,seematti,saree,silk,kancheepuram,handwoven,kottayam,kochi`,
			"name":           "onam 2019 collections",
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)

		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("img_file", filepath.Base(path))
		if err != nil {
			log.Fatal(err)

		}
		_, err = io.Copy(part, file)

		for key, val := range formParams {
			_ = writer.WriteField(key, val)
		}
		err = writer.Close()
		if err != nil {
			log.Fatal(err)

		}
		url := fmt.Sprintf("http://%s/api/upload", serverAddress)

		req, err := http.NewRequest("POST", url, body)
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", writer.FormDataContentType())
		httpClient := http.Client{
			Timeout: 5 * time.Second,
		}
		res, err := httpClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		defer res.Body.Close()
		// Check the response
		Expect(res.StatusCode).To(Equal(200))
	})
})
