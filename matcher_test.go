package scraper

import "testing"

var html = `
<html>
   <body>
       <a href="/link1">Link1</a>
       <h1>Page</h1>
       <a href="/link2">Link2</a>
   </body>
</html>
`

func TestMatcher(t *testing.T) {
	links, err := Match(html)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Errorf("have %d, want 2 links", len(links))
	}
}
