package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHandler struct{}

var page = `
<html>
<body>
  <a href="/page1">p1</a>
  <a href="/page2">p2</a>
  <h1>Page</h1>
</bidy>
</html>
`
var pagedef = `
<html>
<body></body>
</html>
`

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; p {
	case "/":
		w.Write([]byte(page))
		break
	case "/page1":
		w.Write([]byte(pagedef))
	case "/page2":
		w.Write([]byte(pagedef))
	}
}

func createServer() string {
	out := make(chan string)
	go func() {
		ts := httptest.NewServer(testHandler{})
		out <- ts.URL

	}()
	return <-out
}

func TestBrowser(t *testing.T) {
	t.Skip()
	host := createServer()
	browser, err := NewBrowser(host)
	browser.Visit("/page1", func(p Page) {

	})
	if err != nil {
		t.Fatal(err)
	}
	if err := browser.Run(); err != nil {
		t.Fatal(err)
	}
}
