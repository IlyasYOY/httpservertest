package httpservertest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
)

// Matcher simple wrapper wo implement request matching.
type Matcher func(*http.Request) bool

// And is a conjunction for matchers.
func (m Matcher) And(other Matcher) Matcher {
	return func(r *http.Request) bool {
		return m(r) && other(r)
	}
}

// OnGet matches uri with GET mathod.
func OnGet(uri string) Matcher {
	return OnMethod(http.MethodGet, uri)
}

// OnPost matches uri with POST mathod.
func OnPost(uri string) Matcher {
	return OnMethod(http.MethodPost, uri)
}

// OnMethod mathes uri with arbitrary method & uri.
func OnMethod(method, uri string) Matcher {
	return func(r *http.Request) bool {
		return r.URL.RequestURI() == uri && r.Method == method
	}
}

// OnBody matches body of the request.
func OnBody(reader io.Reader) Matcher {
	return func(r *http.Request) bool {
		expected, expectedErr := io.ReadAll(reader)
		if expectedErr != nil {
			panic(fmt.Sprintf("error reading expected body: %v", expectedErr))
		}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			panic(fmt.Sprintf("error reading body: %v", bodyErr))
		}

		return bytes.Equal(expected, body)
	}
}

// OnHeader matches one header. Does not check if it's the only header.
func OnHeader(name, value string) Matcher {
	return func(r *http.Request) bool {
		header, ok := r.Header[name]
		if !ok {
			return false
		}

		return slices.Contains(header, value)
	}
}
