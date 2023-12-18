package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yizeng/gab/chi/crud-gorm/internal/web"

	"github.com/stretchr/testify/assert"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *web.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		assert.Equal(t, actual, expected)
	}
}
