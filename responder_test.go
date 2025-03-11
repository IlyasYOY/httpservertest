package httpservertest_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/IlyasYOY/httpservertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseStatus_200(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(200),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusOK, gotResp.StatusCode)
}

func TestResponseStatus_404(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(404),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, gotResp.StatusCode)
}

func TestResponseBody_404(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(404).
				With(httpservertest.ResponseBody("not found")),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	bytes, readErr := io.ReadAll(gotResp.Body)
	require.NoError(t, readErr)
	assert.EqualValues(t, "not found", bytes)
}

func TestResponseBody_200(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseStatus(200).
				With(httpservertest.ResponseBody("ok")),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	bytes, readErr := io.ReadAll(gotResp.Body)
	require.NoError(t, readErr)
	assert.EqualValues(t, "ok", bytes)
}

func TestResponseStatus_AfterBodyIsUnusual(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseBody("not found").
				With(httpservertest.ResponseStatus(404)),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusOK, gotResp.StatusCode)
}

func TestResponseHeader_ContentType(t *testing.T) {
	server := httpservertest.Start(t)

	server.
		Stub(
			httpservertest.OnGet("/test-url"),
			httpservertest.ResponseBody("not found").
				With(httpservertest.ResponseHeader("Content-Type", "application/json")),
		)

	gotResp, gotErr := http.Get(server.Resolve("/test-url"))
	require.NoError(t, gotErr)
	defer gotResp.Body.Close()

	assert.Equal(t, http.StatusOK, gotResp.StatusCode)
	require.Contains(t, gotResp.Header, "Content-Type")
	header := gotResp.Header["Content-Type"]
	require.Len(t, header, 1)
}
