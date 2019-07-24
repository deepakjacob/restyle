package system_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/deepakjacob/restyle/config"
	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/service"
	"github.com/deepakjacob/restyle/storage"
	"github.com/deepakjacob/restyle/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
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
			log.Fatal(err)
			log.Fatal(err)

		}

		req, err := http.NewRequest("POST", "/upload", body)
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		ctx := config.BootstrapCtx(context.Background())

		logger.Log.Info("init connections to firestore")
		fsClient, err := db.New(ctx)
		if err != nil {
			logger.Log.Fatal("firestore", zap.Error(err))
			return
		}
		logger.Log.Info("init connections to cloud storage")
		csClient, err := storage.New(ctx)
		if err != nil {
			logger.Log.Fatal("cloud storage", zap.Error(err))
			return
		}

		firestoreService := &service.FireStoreServiceImpl{fsClient}
		cloudStorageService := &service.CloudStorageServiceImpl{csClient}
		var mockUser = func(ctx context.Context) (*domain.User, error) {
			return &domain.User{
				Email:  "test@test.com",
				UserID: "090808989898",
			}, nil
		}
		uploadService := &service.UploadServiceImpl{
			FireStoreService:    firestoreService,
			CloudStorageService: cloudStorageService,
			RandStr:             util.RandStr,
			User:                mockUser,
		}

		upload := &handlers.Upload{UploadService: uploadService}
		logger.Log.Debug("test::upload")
		handler := http.HandlerFunc(upload.Handle)
		handler.ServeHTTP(rr, req)

		Expect(err).NotTo(HaveOccurred())
		// Check the response
	})
})
