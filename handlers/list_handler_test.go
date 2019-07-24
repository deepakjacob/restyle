package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/mocks"
	"github.com/deepakjacob/restyle/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("List handler", func() {
	if err := logger.Init(-1, ""); err != nil {
		log.Fatal("logger initialization failed for tests")
	}
	var handler *handlers.List
	var mockListService *service.ListServiceImpl
	var mockFireStoreService *mocks.FireStoreService
	BeforeEach(func() {
		mockFireStoreService = &mocks.FireStoreService{}
		mockListService = &service.ListServiceImpl{
			FireStoreService: mockFireStoreService,
		}
		handler = &handlers.List{
			ListService: mockListService,
		}
	})

	It("/list should responds with code 200", func() {
		data := map[string]string{
			"obj_type":       "saree",
			"material":       "silk",
			"speciality":     "kancheepuram",
			"age_min":        "18",
			"age_max":        "100",
			"name":           "onam 2019 collections",
			"dress_category": "Women",
			"tags":           `2019,onam,seematti,saree,silk,kancheepuram,handwoven,kochi`,
		}

		form := url.Values{}
		for key, val := range data {
			logger.Log.Debug("search::form", zap.String(key, val))
			form.Add(key, val)
		}

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/list", strings.NewReader(form.Encode()))
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			logger.Log.Fatal("post to list images failed", zap.Error(err))
		}

		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))
		Expect(mockFireStoreService.ListCall.Receives.ImgSearch).To(
			Equal(&domain.ImgSearch{
				ObjType: "saree",
			}))
	})

	Context("user lookup fails", func() {
		It("/list should throw internal server error if no user in context", func() {
		})
	})

})
