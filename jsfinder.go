package main

import (
	"os"
)

func main() {

	urlsFile, err := os.Open("urllist.txt") //urlist.txt dosyasını oku
	if err != nil {
		panic(err)
	}
	defer urlsFile.Close()

}
