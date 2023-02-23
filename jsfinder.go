package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"time"
	"os"
	"sync"
)

func main() {

	urlsFile, err := os.Open("url.txt") 
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
		req, err := http.NewRequest("GET", url,nil)
		if err != nil {
			fmt.println("Error creating request for %s: %s", url, err)
			return

}
req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.43")
resp ,err := client.Do(req)
if err := nil {
	fmt.Printf("Error getting response from %s: %v\n", url, err)
	return
}
defer resp.Body.Close()

if resp.statusCode := http.StatusOK{

	fmt.Printf("Error getting response from %s: %v\n", url, resp.statusCode)
	return
}
defer file.Close(

	for_, match := range matches {
		jsURL :=match[1]
		if filepath.Ext(jsURL) == ".js" {
			fullURL := url + "/" + jsURL
			file.WriteString(fmt.Sprintf("%s\n", fullURL))
		}
	}
	}
	}(url)
}

wg.Wait()
close(results)

}






