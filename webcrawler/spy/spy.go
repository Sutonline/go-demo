package spy

import (
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"compress/flate"
	"hash"
	"./impl"
	"time"
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
func Spy() {
	/*defer func() {
		if r := recover(); r != nil {
			log.Println("[E] recover", r)
		}
	}()*/

	// 组装成第一次的url
	client := &http.Client{}
	// ygdy
	ygdy := spyimpl.Ygdy{
		BaseUrl: "http://www.ygdy8.com/html/gndy/dyzz/list_23_",
		PageUrl: "http://www.ygdy8.com/html/gndy/dyzz/list_23_2.html",
		Client:  client}
	spySite(ygdy)

	// pianha
	pianha := spyimpl.Pianha{
		Client: client}
	spySite(pianha)
}

func writeFile(bytes []byte) {
	ioutil.WriteFile("analyze.txt", bytes, 0644)
}

func spySite(spyInterface SpyInterface) {
	totalPages := spyInterface.GetTotalPage()
	log.Printf("一共%d页", totalPages)
	for i := 1; i <= totalPages; i++ {
		log.Printf("开始抓取%d页\n", i)
		spyInterface.FindCurrentPageMovies(spyInterface.GetRequestUrl(i))
		time.Sleep(2 * time.Second)
	}
}
