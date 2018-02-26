package spy


import (
	iconv "github.com/djimenez/iconv-go"
	"../base"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
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
		// body, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
		body := res.Body
		defer res.Body.Close()
		bodyByte, _ := ioutil.ReadAll(body)
		resStr := string(bodyByte)
		output, _ := iconv.ConvertString(resStr, "gb2312", "utf-8")
		return output
	} else {
		return fmt.Sprintf("Get Request not success, code is %s", url)
	}
}
