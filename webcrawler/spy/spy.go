package spy


import (
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"log"
	"net/http"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"io/ioutil"
)

/**
 1. 设置基本网址
 2. 获取页数
 3. 根据页数不断获取
 */
func Spy(baseUrl string) {
	/*defer func() {
		if r := recover(); r != nil {
			log.Println("[E] recover", r)
		}
	}()*/

	// 组装成第一次的url
	url := baseUrl + "2" + ".html"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer res.Body.Close()
		utfBody, err := iconv.NewReader(res.Body, "gb2312", "utf-8")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			doc, _ := goquery.NewDocumentFromReader(utfBody)
			/*divSelection := doc.Find("div.co_content8")
			divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
				val, _ := s.Attr("href")
				fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
			})*/
			// get total page
			page := getPage(doc)
			fmt.Printf("一共%d页\n", page)

			// find movies
			for i := 2; i <= page; i++ {
				log.Printf("开始抓取第%d页", i)
				FindMovies(baseUrl, i)
				time.Sleep(time.Second * 2)
			}
		}
	}
}


func FindMovies(baseUrl string, page int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("抓取第%d页失败，将5秒之后重新尝试\n", page)
			time.Sleep(time.Second * 5)
			FindMovies(baseUrl, page)
		}
	}()
	// concatenate url and page num
	urlStr := baseUrl + strconv.Itoa(page) + ".html"
	log.Printf("抓取的url是: %s", urlStr)

	res, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if res.StatusCode == 200 {
			defer res.Body.Close()
			utfBody, err := iconv.NewReader(res.Body, "gb2312", "utf-8")
			if err != nil {
				fmt.Println(err.Error())
			} else {
				doc, e := goquery.NewDocumentFromReader(utfBody)
				if e != nil {
					bytes, _ := ioutil.ReadAll(utfBody)
					log.Printf("返回的内容: %s", string(bytes))
				}
				// 下面就可以用doc去获取网页里的结构数据了
				// 比如
				divSelection := doc.Find("div.co_content8")
				divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
					val, _ := s.Attr("href")
					fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
				})
			}
		} else {
			log.Printf("请求url: %s失败, statusCode是%d", urlStr, res.StatusCode)
		}
	}
}

func getPage(document *goquery.Document) int {
	totalPage := -1
	regexp.MustCompile("\\(\\_\\)\\(\\d{3,}\\)\\(\\.html\\)")
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