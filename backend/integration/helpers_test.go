//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const baseURL = "http://localhost:8080"

func newExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

type authTransport struct {
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

func authExpect(t *testing.T, token string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: &authTransport{token: token},
		},
	})
}

func registerUser(t *testing.T, username, password, displayName string) string {
	e := newExpect(t)
	obj := e.POST("/api/users/register").
		WithJSON(map[string]string{
			"username":     username,
			"password":     password,
			"display_name": displayName,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	obj.Value("user").Object().NotContainsKey("password_hash")
	return obj.Value("token").String().Raw()
}

func loginUser(t *testing.T, username, password string) string {
	e := newExpect(t)
	return e.POST("/api/users/login").
		WithJSON(map[string]string{
			"username": username,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("token").String().Raw()
}

func findSpaceByName(t *testing.T, token, name string) string {
	t.Helper()
	ae := authExpect(t, token)
	spaces := ae.GET("/api/spaces").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	for i := 0; i < int(spaces.Length().Raw()); i++ {
		obj := spaces.Value(i).Object()
		if obj.Value("name").Raw() == name {
			return obj.Value("id").String().Raw()
		}
	}
	t.Fatalf("space %q not found", name)
	return ""
}
