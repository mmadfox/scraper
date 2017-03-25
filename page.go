package scraper

import (
	"net/http"
	"net/url"
)

type fakeResponseWriter struct {
	Request  *http.Request
	Response *http.Response
	Addr     *url.URL
	Content  string
}

func (rw fakeResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (rw fakeResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (rw fakeResponseWriter) WriteHeader(c int) {
}

type Page struct {
	fakeResponseWriter
}
