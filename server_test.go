package httpservertest_test

import (
	"net/http"
	"testing"

	"github.com/IlyasYOY/httpservertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart_SpinUpSimpleServerAndGet200(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.Response().
				With(httpservertest.ResponseStatus(200)),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusOK, gotResp.StatusCode)
}

func TestServer_Stub_NotMatched(t *testing.T) {
	server := httpservertest.Start(t)

	absentStub := server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseBody("not found").
				With(httpservertest.ResponseStatus(404)),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url-other"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Zero(t, absentStub())
}

func TestServer_Stub_Matched(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.Response().
				With(httpservertest.ResponseStatus(200)),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, 1, stub())
}

func TestServer_Stub_MatchedOneOfTheStubs(t *testing.T) {
	server := httpservertest.Start(t)

	stub := server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.Response().
				With(httpservertest.ResponseStatus(404)),
		)
	absentStub := server.
		Stub(
			httpservertest.OnGet("/test-url-not-matched"),
			httpservertest.Response().
				With(httpservertest.ResponseStatus(200)),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, 1, stub())
	assert.Zero(t, absentStub())
}
