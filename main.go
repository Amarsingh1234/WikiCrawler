package main

import (
    "fmt"
	"os"
    "time"
    "log"

    "github.com/PuerkitoBio/goquery"
    "github.com/gocolly/colly"
	badger "github.com/dgraph-io/badger"
)


func elapsed(what string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", what, time.Since(start))
    }
}

func main() {

    db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
	var idx int=0;


    c := colly.NewCollector(
        colly.AllowedDomains("en.wikipedia.org"),
    )

    // Callback for when a scraped page contains an article element
    c.OnHTML("article", func(e *colly.HTMLElement) {

        metaTags := e.DOM.ParentsUntil("~").Find("meta")
        metaTags.Each(func(_ int, s *goquery.Selection) {

        })

    })

    c.OnHTML(".mw-headline", func(e *colly.HTMLElement) {
      
        // err = db.Update(func(txn *badger.Txn) error {
        //     err := txn.Set([]byte(e.Text), []byte(e.Text))
        //     return err
        // })
        err = db.View(func(txn *badger.Txn) error {
            item, err := txn.Get([]byte(e.Text))
            fmt.Println(err)
        
            err = item.Value(func(val []byte) error {
              fmt.Printf("%s\n", val)
              return nil
            })
          
            return nil
          })
    })

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		
        // Extract the linked URL from the anchor tag
        link := e.Attr("href")
        // fmt.Println(link)
        // Have our crawler visit the linked URL
        c.Visit(e.Request.AbsoluteURL(link))
    })

    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        RandomDelay: 1 * time.Second,
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL.String())
		idx++;
		if idx==1 {
			defer elapsed("page")()
		} else if idx==10 {
			defer elapsed("page")()
		} else if(idx==100){
			defer elapsed("page")()
		} else if(idx==1000){
			defer elapsed("page")()
		} else if idx>1000 {
			os.Exit(3)
		}
    })


    c.Visit("https://en.wikipedia.org")
}