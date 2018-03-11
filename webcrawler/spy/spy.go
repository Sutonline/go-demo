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
	var spy_interface Spy_interface
	ygdy := spyimpl.Ygdy{
		"http://www.ygdy8.com/html/gndy/dyzz/list_23_",
		"http://www.ygdy8.com/html/gndy/dyzz/list_23_2.html",
		"",
		client}
	spy_interface = ygdy
	totalPages := spy_interface.GetTotalPage()
	log.Printf("一共%d页", totalPages)
	for i := 1; i <= totalPages; i++ {
		spy_interface.FindCurrentPageMovies(spy_interface.GetRequestUrl(i))
		time.Sleep(2 * time.Second)
	}

}

func writeFile(bytes []byte) {
	ioutil.WriteFile("analyze.txt", bytes, 0644);
}
