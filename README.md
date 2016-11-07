# scraper
A fast and powerful the scraper for html web pages

## Installation
 $ go get github.com/mmadfox/scraper
 
## Examples

```Go
package main
                               
import (
        "github.com/mmadfox/scraper"    
        "log"                  
        "net/http"
) 
  
func main() {
        log.Println("Imdb scraper")     
        var wc scraper.WorkerCount = 20 
        s, err := scraper.New("http://www.imdb.com/trailers/", wc)
        if err != nil {        
                panic(err)     
        }
        p := "/title/tt{id:[0-9]+}/"    
        s.Mux().HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
                ctx := rw.(*scraper.Context)    
                title := ctx.Doc.Find("h1[itemprop=name]").Text()
                log.Println(title)              
        })
        s.Start()
        s.Block()
}
```
