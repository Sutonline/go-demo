package main

import (
	"./spy"
	"fmt"
)

func main() {
	//spy.ExampleScrape()
	//spy.FindMovies("http://www.ygdy8.com/html/gndy/dyzz/list_23_2.html")
	page := spy.ParsePage("list_23_171.html")
	fmt.Print(page)
}
