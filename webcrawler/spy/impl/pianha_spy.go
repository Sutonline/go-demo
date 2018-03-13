package spyimpl

import (
	"net/http"
	"strconv"
	"fmt"
	"compress/gzip"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"log"
	"io/ioutil"
)

var baseUrl = "http://www.diaomie.com/tv-list-id-1-pg-" //3-order--by--year--letter--area--lang-.html
// var pageUrl = "http://www.diaomie.com/tv-list-id-1-pg-1-order--by--year--letter--area--lang-.html"

type Pianha struct {
	Client *http.Client
}

func (pianha Pianha) GetRequestUrl(page int) string {
	return baseUrl + strconv.Itoa(page) + "-order--by--year--letter--area--lang-.html"
}

func (pianha Pianha) buildRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	req.Header.Set("Host", "www.ygdy8.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests","1")
	req.Header.Set("Cookie", "UM_distinctid=161dc30cac0568-000fb5f0ab92ee-3e3d5100-1fa400-161dc30cac1442; CNZZDATA5783118=cnzz_eid%3D1361120091-1519813985-http%253A%252F%252Fwww.ygdy8.com%252F%26ntime%3D1519813985")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	return req
}

func (pianha Pianha) getResponse(url string) string {
	res, err := pianha.Client.Do(pianha.buildRequest(url))

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if res.StatusCode == 200 {
			gzipReader, _ := gzip.NewReader(res.Body)
			bytes, _ := ioutil.ReadAll(gzipReader)
			return string(bytes)
		}
	}

	return ""
}

// fixed as 800 and judge current page contains movie
func (pianha Pianha) GetTotalPage() int {
	// get response
	return 800
}


func (pianha Pianha) FindCurrentPageMovies(url string) {
	log.Printf("抓取的url是: %s", url)

	responseStr := pianha.getResponse(url)
	// convert to goquery document
	document, e := goquery.NewDocumentFromReader(strings.NewReader(responseStr))
	if e != nil {
		log.Println("build document error", e)
	}

	// find movies
	divSelection := document.Find("ul.show-list")
	playText := divSelection.Find("div.play-txt").Find("h5")
	playText.Find("a").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
	})
}
