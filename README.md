# httpservertest

This is lightweight package to handle HTTP-request stubs.

It's somewhat similar to [WireMock](https://wiremock.org/) , but much less flexible.

Simple usage:

```go
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
```

More examples in tests.
