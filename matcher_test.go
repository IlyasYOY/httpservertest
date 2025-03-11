package httpservertest_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/IlyasYOY/httpservertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnGet_Mathed(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.NotZero(t, stub())
}

func TestOnGet_NotMatched(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Post(
		server.Resolve("/test-url"),
		"application/json",
		strings.NewReader(""),
	)
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusOK, gotResp.StatusCode)
	assert.Zero(t, stub())
}

func TestOnPost_Matched(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnPost("/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Post(
		server.Resolve("/test-url"),
		"application/json",
		strings.NewReader(""),
	)
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.NotZero(t, stub())
}

func TestOnPost_NotMatchedd(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnPost("/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Zero(t, stub())
}

func TestOnMethod_Matched(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnMethod("GET", "/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.NotZero(t, stub())
}

func TestMatcher_With(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnMethod("POST", "/test-url").
				And(httpservertest.OnBody(strings.NewReader("test-body"))).
				And(httpservertest.OnHeader("Content-Type", "application/json")),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Post(
		server.Resolve("/test-url"),
		"application/json",
		strings.NewReader("test-body"),
	)
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.NotZero(t, stub())
}
