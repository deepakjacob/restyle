package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/mocks"
	"github.com/deepakjacob/restyle/service"
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
		form := map[string]string{
			"obj_type": "saree",
		}
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/list", form, "img_file", path)
		if err != nil {
			logger.Log.Fatal("listing images failed", zap.Error(err))
		}

		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(200))
		Expect(mockFireStoreService.ListCall.Receives.SearchAttrs).To(Equal(mapImgAttrs(form)))
	})

	Context("user lookup fails", func() {
		It("/list should throw internal server error if no user in context", func() {
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
