package main

import (
	"github.com/PuerkitoBio/goquery"
)

type DummyParser struct {
}

func (d DummyParser) ParsePage(doc *goquery.Document) ScrapeResult {
	data := ScrapeResult{}
	data.Title = doc.Find("title").First().Text()
	data.H1 = doc.Find("h1").First().Text()
	return ScrapeResult{}
}

func main() {
	d := DummyParser{}
	Crawl("https://www.daum.net", d, 1)
}
