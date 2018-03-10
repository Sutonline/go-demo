package spy

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"runtime/debug"
	"io/ioutil"
	"strings"
	"compress/gzip"
	"os/exec"
	bytes2 "bytes"
	"io"
	"compress/flate"
	"hash"
	"./impl"
)

type Reader struct {
	r            flate.Reader
	decompressor io.ReadCloser
	digest       hash.Hash32
	size         uint32
	flg          byte
	buf          [512]byte
	err          error
	multistream  bool
}

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
	var spy_interface Spy_interface
	ygdy := spyimpl.Ygdy{
		"http://www.ygdy8.com/html/gndy/dyzz/list_23_",
		"http://www.ygdy8.com/html/gndy/dyzz/list_23_2.html",
		"",
		client}
	spy_interface = ygdy
	spy_interface.GetTotalPage()
	url := baseUrl + "1" + ".html"

	res := getResponse(client, url)

	if res == nil {
		log.Printf("get [%s] response is nil\n", url)
	} else {
		defer res.Body.Close()
		gzipReader, _ := gzip.NewReader(res.Body)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(convertReader(gzipReader, "gb2312", "utf-8")))
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

	res := getResponse(client, urlStr)

	if res == nil {
		log.Printf("get [%s] response is nil\n", urlStr)
	} else {
		if res.StatusCode == 200 {
			defer res.Body.Close()
			gzipReader, _ := gzip.NewReader(res.Body)
			doc, e := goquery.NewDocumentFromReader(strings.NewReader(convertReader(gzipReader, "gb2312", "utf-8")))
			if e != nil {
				debug.PrintStack()
				log.Fatal("goquery build document ", e)
			}

			divSelection := doc.Find("div.co_content8")
			divSelection.Find("a.ulink").Each(func(i int, s *goquery.Selection) {
				val, _ := s.Attr("href")
				fmt.Printf("find link name:%s, href: %s\n", s.Text(), val)
			})
		} else {
			log.Printf("请求url: %s失败, statusCode是%d", urlStr, res.StatusCode)
		}
	}
}

func getResponse(client *http.Client, urlStr string) *http.Response {
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	req.Header.Set("Host", "www.ygdy8.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cookie", "UM_distinctid=161dc30cac0568-000fb5f0ab92ee-3e3d5100-1fa400-161dc30cac1442; CNZZDATA5783118=cnzz_eid%3D1361120091-1519813985-http%253A%252F%252Fwww.ygdy8.com%252F%26ntime%3D1519813985")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if res.StatusCode == 200 {
			return res
		}
	}

	return nil
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

func convertReader(reader io.Reader, fromCharset string, toCharset string) string {
	bytes, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Fatal("reading bytes", e)
	}
	//writeFile(bytes)
	command := exec.Command("iconv", "-f", "gb2312", "-t", "utf-8", "-c")
	command.Stdin = strings.NewReader(string(bytes))
	var out bytes2.Buffer
	command.Stdout = &out
	command.Run()

	return string(out.Bytes())
}

func writeFile(bytes []byte) {
	ioutil.WriteFile("analyze.txt", bytes, 0644);
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
