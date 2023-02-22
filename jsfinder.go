package main

import (
	"os"
	"sync"
)

func main() {

	urlsFile, err := os.Open("url.txt") //urlist.txt dosyasını oku
	if err != nil {
		panic(err)
	}
	defer urlsFile.Close()

	results := make(chan string)

	sem := make(chan struct{}, 50)

	var wg sync.WaitGroup

}
