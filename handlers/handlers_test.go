package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/mocks"
	"github.com/deepakjacob/restyle/oauth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("Auth handler", func() {
	var handler *handlers.OAuth2
	var mockProvider *oauth.Provider
	var mockUserService *mocks.UserService
	var mockConfig *mocks.Config
	var mockGoogleClient *mocks.GoogleClient
	BeforeEach(func() {
		mockGoogleClient = &mocks.GoogleClient{}
		mockGoogleClient.GetCall.Returns.GoogleUser = &oauth.GoogleUser{
			Email: "test@test.com",
		}
		mockConfig = &mocks.Config{}
		mockConfig.ExchangeCall.Returns.Token = &oauth2.Token{
			AccessToken: "access_token",
		}
		mockProvider = &oauth.Provider{
			Config:     mockConfig,
			HTTPClient: mockGoogleClient,
		}
		mockUserService = &mocks.UserService{}
		mockUserService.FindCall.Returns.User = &domain.User{
			Email: "test@test.com",
		}
		handler = &handlers.OAuth2{
			Provider:    mockProvider,
			RandStr:     mocks.RandStr,
			UserService: mockUserService,
		}
	})

	It("responds with code 307 and string redirect url", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth", nil)
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(307))
		request := &http.Request{
			Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]},
		}
		cookie, _ := request.Cookie("state")
		Expect(cookie).Should(BeAssignableToTypeOf(&http.Cookie{}))
		Expect(("a_random_string")).To(Equal(cookie.Value))
		// TODO: fix this
		// Expect(cookie.Domain).To(Equal("/"))
		// Expect(cookie.HttpOnly).Should(BeTrue())
		Expect(mockConfig.AuthCodeURLCall.Receives.State).To(Equal("a_random_string"))
	})
	It("callback handler fetches google user", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth/callback?state=userstate&code=usercode", nil)
		handler := http.HandlerFunc(handler.HandleCallback)
		req.AddCookie(&http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Value:    "userstate",
			Name:     "state",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
		})
		handler.ServeHTTP(resp, req)
		Expect(mockConfig.ExchangeCall.Receives.Code).To(Equal("usercode"))
		Expect(mockGoogleClient.GetCall.Receives.URL).To(
			Equal("https://www.googleapis.com/oauth2/v2/userinfo?access_token=access_token"))
		Expect(mockUserService.FindCall.Receives.Email).To(
			Equal(mockGoogleClient.GetCall.Returns.GoogleUser.Email))
		Expect(resp.Code).To(Equal(307))
	})
})
