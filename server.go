package httpservertest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Start [httptest.NewServer] tied to [testing.T] life-cycle.
// Later it might be used to validate calls to it.
func Start(t *testing.T) *server {
	t.Helper()
	handler := &matchingHandlers{}
	testServer := httptest.NewServer(handler)
	t.Cleanup(func() {
		testServer.Close()
	})
	return &server{
		t:          t,
		testServer: testServer,
		handlers:   handler,
	}
}

type server struct {
	t          *testing.T
	testServer *httptest.Server
	handlers   *matchingHandlers
}

// Resolve an URL of the path relative to the bootstraped server.
// Returns: http://localhost:port/path
// This function adds path as with without any forward slash processing.
func (s *server) Resolve(path string) string {
	return s.testServer.URL + path
}

type matchingHandlers struct {
	hs []*matchingHandler
}

// ServeHTTP implements http.Handler.
func (m *matchingHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range m.hs {
		h.ServeHTTP(w, r)
	}
}

type matchingHandler struct {
	matcher   Matcher
	responder Responder
	matched   int
}

// ServeHTTP implements http.Handler.
func (m *matchingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.matcher(r) {
		m.responder(w)
		m.matched++
	}
}

type stub func() int

// Stub a call with a given matcher & responder.
func (s *server) Stub(m Matcher, r Responder) stub {
	mh := &matchingHandler{
		matcher:   m,
		responder: r,
		matched:   0,
	}
	s.handlers.hs = append(s.handlers.hs, mh)
	return func() int {
		return mh.matched
	}
}
