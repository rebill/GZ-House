package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var area = "tianhe"
var maxPage = 50

func main() {
	fName := fmt.Sprintf("gz_house_%s.csv", area)
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Title", "District", "Bizcircle", "Tag", "Built", "TotalPrice"})

	// Instantiate default collector
	c := colly.NewCollector(
	// Visit only domains
	// colly.AllowedDomains("gz.lianjia.com"),
	// colly.Async(true),
	)

	c.SetRequestTimeout(30 * time.Second)

	// Limit the number of threads started by colly to two
	// c.Limit(&colly.LimitRule{
	// Parallelism: 2,
	// RandomDelay: 5 * time.Second,
	// })

	extensions.RandomUserAgent(c)
	extensions.Referrer(c)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Retry()
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("li.xiaoquListItem", func(e *colly.HTMLElement) {

		title := e.ChildText(".title a")
		// fmt.Println("title:", title)

		district := e.ChildText(".district")
		// fmt.Println("district:", district)

		bizcircle := e.ChildText(".bizcircle")
		// fmt.Println("bizcircle:", bizcircle)

		positionInfo := e.ChildText(".positionInfo")
		r := strings.NewReplacer("未知", "1970", "年建成", "")
		built := strings.TrimSpace(r.Replace(strings.Split(positionInfo, "/")[1]))
		// fmt.Println("built:", built)

		tag := e.ChildText(".tagList span")
		// fmt.Println("tag:", tag)

		totalPrice := e.ChildText(".totalPrice span")
		totalPrice = strings.Replace(totalPrice, "暂无", "0", -1)
		// fmt.Println("totalPrice:", totalPrice)

		// fmt.Println()

		writer.Write([]string{
			title,
			district,
			bizcircle,
			tag,
			built,
			totalPrice,
		})

	})

	urlFmt := "https://gz.lianjia.com/xiaoqu/%s/pg%dcro11/"
	for i := 1; i <= maxPage; i++ {
		c.Visit(fmt.Sprintf(urlFmt, area, i))
		time.Sleep(time.Duration(rand.Intn(int(5 * time.Second))))
	}

	// Wait until threads are finished
	c.Wait()

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
