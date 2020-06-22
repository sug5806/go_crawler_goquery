package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ScrapeResult struct {
	URL   string
	Title string
	H1    string
}

type Parser interface {
	ParsePage(*goquery.Document) ScrapeResult
}

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func extractLinks(doc *goquery.Document) []string {
	var foundUrls []string
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	}
	return foundUrls
}

func resolveRelative(baseURL string, hrefs []string) []string {
	var internalUrls []string

	for _, href := range hrefs {
		if strings.HasPrefix(href, baseURL) {
			internalUrls = append(internalUrls, href)
		}

		if strings.HasPrefix(href, "/") {
			resolvedURL := fmt.Sprintf("%s%s", baseURL, href)
			internalUrls = append(internalUrls, resolvedURL)
		}
	}

	return internalUrls
}

//func crawlPage(baseURL, targetURL string, parser Parser, token chan struct{}) ([]string, ScrapeResult) {
//	token <- struct{}{}
//	fmt.Println("Requesting: ", targetURL)
//	resp, _ := getRequest(targetURL)
//	<-token
//
//	log.Println("baseUrl : ", baseURL)
//	log.Println("targetURL : ", targetURL)
//
//	doc, _ := goquery.NewDocumentFromReader(resp.Body)
//	pageResults := parser.ParsePage(doc)
//	links := extractLinks(doc)
//	foundUrls := resolveRelative(baseURL, links)
//
//	return foundUrls, pageResults
//
//}

func crawlPage(baseURL, targetURL string, parser Parser) ([]string, ScrapeResult) {
	//token <- struct{}{}
	log.Println("Requesting: ", targetURL)
	resp, _ := getRequest(targetURL)
	//<-token
	fmt.Println()
	log.Println("baseUrl : ", baseURL)
	log.Println("targetURL : ", targetURL)

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	pageResults := parser.ParsePage(doc)
	links := extractLinks(doc)
	foundUrls := resolveRelative(baseURL, links)

	return foundUrls, pageResults

}

func parseStartURL(u string) string {
	parsed, _ := url.Parse(u)
	log.Println("parsed.Scheme : ", parsed.Scheme)
	log.Println("parsed.Host : ", parsed.Host)
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

func Crawl(startURL string, parser Parser, concurrency int) []ScrapeResult {
	var results []ScrapeResult
	//workList := make(chan []string)
	//var n int
	//n++
	//var tokens = make(chan struct{}, concurrency)
	//go func() { workList <- []string{startURL}}()
	//fmt.Println("workList : ", &workList)
	//seen := make(map[string]bool)
	baseDomain := parseStartURL(startURL)
	foundLinks, pageResults := crawlPage(baseDomain, startURL, parser)
	results = append(results, pageResults)

	if foundLinks != nil {
		log.Println("foundLInks : ", foundLinks)
	}


	//for ; n >0; n-- {
	//	list := <-workList
	//	for _, link := range list {
	//		if !seen[link] {
	//			seen[link] = true
	//			n++
	//			go func(baseDomain, link string, parser Parser, token chan struct{}) {
	//				foundLinks, pageResults := crawlPage(baseDomain, link, parser, token)
	//				results = append(results, pageResults)
	//				if foundLinks != nil {
	//					workList <- foundLinks
	//				}
	//			}(baseDomain, link, parser, tokens)
	//		}
	//	}
	//}
	return results
}

