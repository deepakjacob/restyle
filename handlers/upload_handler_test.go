package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Upload handler", func() {
	if err := logger.Init(-1, ""); err != nil {
		log.Fatal("logger initialization failed for tests")
	}
	var handler *handlers.Upload
	var mockUploadService *mocks.UploadService
	BeforeEach(func() {
		mockUploadService = &mocks.UploadService{}
		handler = &handlers.Upload{
			UploadService: mockUploadService,
		}
	})

	It("/upload should responds with code 200", func() {
		path, _ := os.Getwd()
		path += "/test_resources/img_1.jpg"
		form := map[string]string{
			"obj_type":       "saree",
			"material":       "silk",
			"speciality":     "kancheepuram",
			"age_min":        "18",
			"age_max":        "100",
			"name":           "onam 2019 collections",
			"dress_category": "Women",
			"tags":           `2019,onam,seematti,saree,silk,kancheepuram,handwoven,kochi`,
		}
		resp := httptest.NewRecorder()
		req, err := uploadRequest("/upload", form, "img_file", path)
		if err != nil {
			logger.Log.Fatal("upload failed", zap.Error(err))
		}
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))
		Expect(mockUploadService.UploadCall.Receives.ImgAttrs).To(Equal(mapImgAttrs(form)))
	})

	Context("user lookup fails", func() {
		BeforeEach(func() {
			mockUploadService.UploadCall.Returns.Error = errors.New("upload error")
		})

		It("/upload should throw internal server error on upload failure", func() {
		})
	})

})

func mapImgAttrs(m map[string]string) *domain.ImgAttrs {
	return &domain.ImgAttrs{
		ObjType:       m["obj_type"],
		Material:      m["material"],
		Speciality:    m["speciality"],
		AgeMin:        18,
		AgeMax:        100,
		Name:          m["name"],
		DressCategory: m["dress_category"],
		Tags:          strings.Split(m["tags"], ","),
	}
}

func uploadRequest(uri string, form map[string]string, field, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range form {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
