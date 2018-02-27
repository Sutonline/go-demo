package spy


import (
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"../base"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
)

func Spy(url string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[E]", r)
		}
	}()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", base.GetRandomUserAgent())
	client := http.DefaultClient
	res, e := client.Do(req)
	if e != nil {
		fmt.Printf("Get request %s failed: %s", url, e)
		return "request failed"
	}

	if res.StatusCode == 200 {
		body := res.Body
		defer res.Body.Close()
		// utfBody, _ := iconv.NewReader(body, "gbk2312", "utf-8")
		bytes, _ := ioutil.ReadAll(body)
		bodyStr := string(bytes)
		output, _ := iconv.ConvertString(bodyStr, "gb2312", "utf-8")

		// find movies
		reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			fmt.Printf("goquery build document from read occur error: %s", err)
			return ""
		}
		doc.Find("div.co_content8").Find("a.ulink").Each(func (i int, s *goquery.Selection) {
			val, exists := s.Attr("href")
			if exists {
				fmt.Printf("find movie name: %s, href: %s", s.Text(), val)
			}
		})

		return output
	} else {
		return fmt.Sprintf("Get Request not success, code is %s", url)
	}
}


func FindMovies(urlStr string) {
	res, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer res.Body.Close()
		utfBody, err := iconv.NewReader(res.Body, "gb2312", "utf-8")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			doc, _ := goquery.NewDocumentFromReader(utfBody)
			// 下面就可以用doc去获取网页里的结构数据了
			// 比如
			divSelection := doc.Find("div.co_content8")
			divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
					val, _ := s.Attr("href")
					fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
				})
		}
	}
}

func getPage(document goquery.Document) int {
	totalPage := -1
	regexp.MustCompile("\\(\\_\\)\\(\\d{3,}\\)\\(\\.html\\)")
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "末页" {
			val, exists := s.Attr("href")
			if !exists {
				return
			}
			fmt.Print(val)
		}
	})
	return totalPage
}

func ParsePage(pageStr string) int {
	compile := regexp.MustCompile("(_)(\\d{3,})(.html)")
	submatch := compile.FindStringSubmatch(pageStr)
	for i, v := range submatch {
		fmt.Printf("index %d find %s\n", i, v)
	}
	return -1
}