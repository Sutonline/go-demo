package main

import (
	"./spy"
	"io/ioutil"
)

func main() {
	s := spy.Spy("http://www.dytt8.net")
	ioutil.WriteFile("index.txt", []byte(s), 0644)
}
