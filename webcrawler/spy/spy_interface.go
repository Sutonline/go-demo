package spy

/**
  1. 获取总页数
  2. 根据总页数生成各个页的url
  3. 获取当前页的电影
 */
type Spy_interface interface {

	// 根据baseUrl获取进行请求的url
	GetRequestUrl(page int) string

	// 使用url构造request对象
	//BuildRequest(url string) *http.Request

	// 进行转换字符集等操作获取正确的response
	// GetResponse(url string) string

	// 获取总页数
	GetTotalPage() int

	// 将当前页的url的movie解析出来
	FindCurrentPageMovies(url string)

}
