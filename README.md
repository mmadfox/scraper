# scraper
A fast and powerful the scraper for html web pages

## Installation
 $ go get github.com/mmadfox/scraper
 
## Examples
### https://trakt.tv/
```Go
package main
                               
import (
        "github.com/mmadfox/scraper"    
        "log"                  
        "net/http"             
        "time"
) 
  
func main() {
        log.Println("Trakt Tv scraper") 
        var wc scraper.WorkerCount = 5  
        h := http.Header{}     
        h.Add("Referer", "https://trakt.tv")
  
        s, err := scraper.New("https://trakt.tv/movies/trending", wc)
        if err != nil {
                panic(err)
        }
        p := `/movies/{movieName:(.*)\-[0-9]+}`
        s.Mux().HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
                ctx := rw.(*scraper.Context)    
                log.Println("Got the url", ctx.Addr.String())
        })
        s.SetHeader(h)
        go func() {
                for {
                        select {                        
                        case <-time.After(1 * time.Minute):
                                s.StopAndClose()                
                                return                          
                        }
                }
        }()
        s.Start()
        s.Block()
        log.Println("Stop")
}

```
###http://www.imdb.com/
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
