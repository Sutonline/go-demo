package spyimpl

import (
	"strconv"
	"net/http"
	"compress/gzip"
	"../../base"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"log"
	"strings"
	"time"
	"runtime/debug"
	"fmt"
)

type Ygdy struct {
	baseUrl string
	pageUrl string
	responseStr string
	client *http.Client
}

func (ygdy Ygdy) GetRequestUrl(page int) string {
	return ygdy.baseUrl + strconv.Itoa(page) + ".html"
}

func (ygdy Ygdy) buildRequest(url string) *http.Request {
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

func (ygdy Ygdy) GetResponse(url string) string {
	res, err := ygdy.client.Do(ygdy.buildRequest(url))

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if res.StatusCode == 200 {
			gzipReader, _ := gzip.NewReader(res.Body)
			ygdy.responseStr = base.ConvertReader(gzipReader, "gb2312", "utf-8")
			return ygdy.responseStr
		}
	}

	return ""
}

func (ygdy Ygdy) GetTotalPage() int {
	resStr := ygdy.GetResponse(ygdy.pageUrl)
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(resStr))
	return getPage(document)
}

func (ygdy Ygdy) FindCurrentPageMovies(url string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[E] recover", r)
			log.Printf("抓取%s失败，将5秒之后重新尝试\n", url)
			time.Sleep(time.Second * 5)
			ygdy.FindCurrentPageMovies(url)
		}
	}()

	log.Printf("抓取的url是: %s", url)

	res := ygdy.GetResponse(url)

	if res == "" {
		log.Printf("get [%s] response is nil\n", url)
	} else {
			doc, e := goquery.NewDocumentFromReader(strings.NewReader(res))
			if e != nil {
				debug.PrintStack()
				log.Fatal("goquery build document ", e)
			}

			divSelection := doc.Find("div.co_content8")
			divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
				val, _ := s.Attr("href")
				fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
			})
		}
	}

func getPage(document *goquery.Document) int {
	totalPage := -1
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "末页" {
			val, exists := s.Attr("href")
			if !exists {
				return
			}
			totalPage = ParsePage(val)
		}
	})
	return totalPage
}

func ParsePage(pageStr string) int {
	compile := regexp.MustCompile("(_)(\\d{3,})(.html)")
	submatch := compile.FindStringSubmatch(pageStr)
	if len(submatch) < 3 {
		return -1
	}
	page, e := strconv.Atoi(submatch[2])
	if e != nil {
		log.Fatal("获取分页的时候发生错误", e)
	}

	return page
}