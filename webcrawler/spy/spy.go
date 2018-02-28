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
	"runtime/debug"
	"io/ioutil"
	"strings"
	"os"
)

var converter, _ = iconv.NewConverter("GB18030", "utf-8")

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
	client := &http.Client{}
	url := baseUrl + "5" + ".html"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	req.Header.Set("Host", "www.ygdy8.com")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer res.Body.Close()
		// reader, _ := gzip.NewReader(res.Body)
		bytes, _ := ioutil.ReadAll(res.Body)
		s := string(bytes)
		utfStr, _ := converter.ConvertString(s)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(utfStr))
			// get total page
			page := getPage(doc)
			fmt.Printf("一共%d页\n", page)

			// find movies
			for i := 1; i <= page; i++ {
				log.Printf("开始抓取第%d页", i)
				FindMovies(client, baseUrl, i)
				time.Sleep(time.Second * 2)
			}
		}
	}
}


func FindMovies(client *http.Client, baseUrl string, page int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[E] recover", r)
			log.Printf("抓取第%d页失败，将5秒之后重新尝试\n", page)
			time.Sleep(time.Second * 5)
			FindMovies(client, baseUrl, page)
		}
	}()
	// concatenate url and page num
	urlStr := baseUrl + strconv.Itoa(page) + ".html"
	log.Printf("抓取的url是: %s", urlStr)

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	req.Header.Set("Host", "www.ygdy8.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Refer", baseUrl + strconv.Itoa(page - 1) + ".html")
	req.Header.Set("Upgrade-Insecure-Requests","1")
	req.Header.Set("Cookie", "UM_distinctid=161dc30cac0568-000fb5f0ab92ee-3e3d5100-1fa400-161dc30cac1442; CNZZDATA5783118=cnzz_eid%3D1361120091-1519813985-http%253A%252F%252Fwww.ygdy8.com%252F%26ntime%3D1519813985")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if res.StatusCode == 200 {
			defer res.Body.Close()
			bytes, _ := ioutil.ReadAll(res.Body)
			utfStr, _ := converter.ConvertString(string(bytes))

			if err != nil {
				fmt.Println(err.Error())
			} else {
				doc, e := goquery.NewDocumentFromReader(strings.NewReader(utfStr))
				if e != nil {
					debug.PrintStack()
					log.Fatal("goquery build document ", e)
				}
				// 下面就可以用doc去获取网页里的结构数据了
				// 比如
				divSelection := doc.Find("div.co_content8")
				divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
					val, _ := s.Attr("href")
					fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
				})

				// 记录文件
				_, err := os.Create("/opt/git/go-demo/webcrawler/html/" + strconv.Itoa(page) + ".txt")
				if err == nil {
					ioutil.WriteFile("/opt/git/go-demo/webcrawler/html/"+ strconv.Itoa(page) + ".txt", []byte(utfStr), 0644)
				} else {
					log.Print("create file error", err)
				}
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