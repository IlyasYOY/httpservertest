package httpservertest

import "net/http"

// Responder is a simple way to define response. Just a wrapper around a
// function.
type Responder func(http.ResponseWriter)

// With allows you to combine responses as you wish.
func (r Responder) With(other Responder) Responder {
	return func(w http.ResponseWriter) {
		r(w)
		other(w)
	}
}

// Response empty resonse.
func Response() Responder {
	return func(w http.ResponseWriter) {}
}

// ResponseBody specifies body of the response as a string.
func ResponseBody(body string) Responder {
	return func(w http.ResponseWriter) {
		_, err := w.Write([]byte(body))
		if err != nil {
			panic(err)
		}
	}
}

// ResponseStatus specifies status code.
// this method must be called first due to the [http.ResponseWriter]. Otherwise
// we will always get 200.
func ResponseStatus(statusCode int) Responder {
	return func(w http.ResponseWriter) {
		w.WriteHeader(statusCode)
	}
}

// ResponseHeader specifies header in the reponse.
func ResponseHeader(name, value string) Responder {
	return func(w http.ResponseWriter) {
		w.Header().Add(name, value)
	}
}
