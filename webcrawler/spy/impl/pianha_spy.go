package spyimpl

import (
	"net/http"
	"strconv"
)

var baseUrl = "http://www.diaomie.com/tv-list-id-1-pg-" //3-order--by--year--letter--area--lang-.html
var pageUrl = "http://www.diaomie.com/tv-list-id-1-pg-1-order--by--year--letter--area--lang-.html"

type Pianha struct {
	Client *http.Client
}

func (pianha Pianha) GetRequestUrl(page int) string {
	return baseUrl + strconv.Itoa(page) + "-order--by--year--letter--area--lang-.html"
}

func (pianha Pianha) GetTotalPage() int {
	return -1
}

