package handlers_test

import (
	"errors"
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
	var mockSigner *mocks.Signer

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

	It("responds with code 307 and string redirect url", func() {
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

	It("callback handler set a cookie named _ut and redirects", func() {
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
		//   Expect(resp.Code).To(Equal(307))
	})

	It("should return http forbidden if state token is missing in url", func() {
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

	It("should return http forbidden if code is missing in url", func() {
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

		It("should throw internal server error", func() {
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
			mockSigner.SignEncryptCall.Returns.Err = errors.New("Signing failed for some reason")
		})

		It("should throw internal server error ", func() {
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
