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

scanner :=bufio.NewScanner(urlsFile)
for scanner.Scan(){
	url:=scanner.Text()
	wg.Add(1)

	go func (url string) {

		defer wg.Done()
		sem <- struct{}{}
		defer func(){<-sem}()

		client:= &http.Client{
			timeout: 5 * time.Second,
		}

}



}




