package handlers_test

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/mocks"
	"github.com/deepakjacob/restyle/oauth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("Auth handler", func() {
	if err := logger.Init(-1, ""); err != nil {
		log.Fatal("logger initialization failed for tests")
	}
	var handler *handlers.OAuth2
	var mockConfig *mocks.Config
	var mockSigner *mocks.Signer
	var mockGoogleClient *mocks.GoogleClient
	var mockProvider *oauth.ProviderImpl
	var mockUserService *mocks.UserService
	BeforeEach(func() {
		mockGoogleClient = &mocks.GoogleClient{}
		mockGoogleClient.GetCall.Returns.GoogleUser = &oauth.GoogleUser{
			Email: "test@test.com",
		}
		mockConfig = &mocks.Config{}
		mockConfig.ExchangeCall.Returns.Token = &oauth2.Token{
			AccessToken: "access_token",
		}
		mockProvider = &oauth.ProviderImpl{
			Config:     mockConfig,
			HTTPClient: mockGoogleClient,
		}
		mockUserService = &mocks.UserService{}
		mockUserService.FindCall.Returns.User = &domain.User{
			Email:  "test@test.com",
			UserID: "12345",
		}
		mockSigner = &mocks.Signer{}
		mockSigner.SignEncryptCall.Returns.SignedValue = "signed_user"
		handler = &handlers.OAuth2{
			Provider:    mockProvider,
			RandStr:     mocks.RandStr,
			UserService: mockUserService,
			Signer:      mockSigner,
		}
	})

	It("/auth should responds with code 307 and string redirect url", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth", nil)
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(307))
		request := &http.Request{
			Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]},
		}
		cookie, err := request.Cookie("state")
		Expect(err).ToNot(HaveOccurred())
		Expect(("a_random_string")).To(Equal(cookie.Value))
		Expect(mockConfig.AuthCodeURLCall.Receives.State).To(Equal("a_random_string"))
	})

	It("/auth/callback should set a cookie named _ut and redirects", func() {
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
		Expect(mockSigner.SignEncryptCall.Receives.User.Email).To(Equal("test@test.com"))
		request := &http.Request{
			Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]},
		}

		cookie, err := request.Cookie("_ut")
		Expect(err).ToNot(HaveOccurred())
		Expect(("signed_user")).To(Equal(cookie.Value))
		Expect(resp.Code).To(Equal(307))
	})

	It("/auth/callback should return http forbidden if state token is missing in url", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth/callback?code=usercode", nil)
		handler := http.HandlerFunc(handler.HandleCallback)
		handler.ServeHTTP(resp, req)
		request := &http.Request{
			Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]},
		}
		_, err := request.Cookie("state")
		Expect(err).To(HaveOccurred())
		Expect(resp.Code).To(Equal(403))
	})

	It("/auth/callback should return http forbidden if code is missing in url", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth/callback?state=userstate", nil)
		req.AddCookie(&http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Value:    "userstate",
			Name:     "state",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
		})
		handler := http.HandlerFunc(handler.HandleCallback)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(403))
	})

	Context("user lookup fails", func() {
		BeforeEach(func() {
			mockUserService.FindCall.Returns.Error = errors.New("user not found")
		})

		It("/auth/callback should throw internal server error on user lookup failure", func() {
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
			Expect(resp.Code).To(Equal(500))
		})
	})
	Context("when signing fails", func() {
		BeforeEach(func() {
			mockSigner.SignEncryptCall.Returns.Err = errors.New("signing failed")
		})

		It("/auth/callback should throw internal server error when signing fails", func() {
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
			Expect(resp.Code).To(Equal(500))
		})
	})
})
