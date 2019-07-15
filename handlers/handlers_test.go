package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth handler", func() {
	var handler *handlers.OAuth2
	var mockProvider *mocks.Provider

	BeforeEach(func() {
		mockProvider = &mocks.Provider{
			Config: &oauth2.Config{
				Endpoint:     google.Endpoint,
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
				RedirectURL:  "server_auth_callback",
				Scopes: []string{
					"scope1",
					"scope2",
				},
			},
		}
		handler = &handlers.OAuth2{
			Provider: mockProvider,
			RandStr:  mocks.RandStr,
		}
	})

	It("responds with code 307 and string redirect url", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth", nil)
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(307))
		u, _ := url.Parse(resp.Header().Get("Location"))
		Expect(u.Scheme).To(Equal("https"))
		Expect(u.Host).To(Equal("accounts.google.com"))
		q := u.Query()
		Expect((q["state"])[0]).To(Equal("a_random_string"))
		Expect((q["redirect_uri"])[0]).To(Equal("server_auth_callback"))
	})
	It("sets a matching cookie in response", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth", nil)
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		request := &http.Request{Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]}}
		cookie, _ := request.Cookie("state")
		u, _ := url.Parse(resp.Header().Get("Location"))
		q := u.Query()
		Expect((q["state"])[0]).To(Equal(cookie.Value))
		Expect(cookie.Value).To(Equal("a_random_string"))
	})
	It("should fail with not authorized", func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/auth/callback", nil)
		handler := http.HandlerFunc(handler.Handle)
		handler.ServeHTTP(resp, req)
		request := &http.Request{Header: http.Header{"Cookie": resp.HeaderMap["Set-Cookie"]}}
		cookie, _ := request.Cookie("state")
		u, _ := url.Parse(resp.Header().Get("Location"))
		q := u.Query()
		Expect((q["state"])[0]).To(Equal(cookie.Value))
		Expect(cookie.Value).To(Equal("a_random_string"))
	})

})
