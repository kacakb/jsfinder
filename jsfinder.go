package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {

	var limit int
	flag.IntVar(&limit, "c", 25, "concurrency limit")

	var urlsFilePath string
	flag.StringVar(&urlsFilePath, "l", "", "filename to read URLS from")
	// Open url file

	var silent bool
	flag.BoolVar(&silent, "s", false, "Program is running in silent mode")

	flag.Parse()

	if urlsFilePath == "" {
		fmt.Println("Please provide a file containing URLs wtih the -l flag")
		return
	}
	urlsFile, err := os.Open(urlsFilePath)
	if err != nil {
		panic(err)
	}
	defer urlsFile.Close()

	// Channel  goroutines
	results := make(chan string)
	// Semaphore to limit
	sem := make(chan struct{}, limit)

	//  wait for all goroutines to finish
	var wg sync.WaitGroup

	// Read URLs file line by line
	scanner := bufio.NewScanner(urlsFile)
	for scanner.Scan() {
		url := scanner.Text()

		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			//  semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			// Create HTTP client with custom User-Agent header and timeout
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request for %s: %v\n", url, err)
				return
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")

			// Send GET request to URL
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error getting response from %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			// Check status code for successful request
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error getting response from %s: status code %d\n", url, resp.StatusCode)
				return
			}

			// Read response body
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				if !silent {
					fmt.Printf("Error reading response from %s: %v\n", url, err)
				}

				return
			}
			bodyString := string(bodyBytes)

			// Find JavaScript files using regular expression
			re := regexp.MustCompile(`src="([^"]+\.js)"`)

			//js files usually "src=.js$"

			matches := re.FindAllStringSubmatch(bodyString, -1)
			if len(matches) > 0 {
				// Write JavaScript file names to file
				file, err := os.OpenFile("found.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Error opening file for writing: %v\n", err)
					return
				}
				defer file.Close()

				for _, match := range matches {
					jsURL := match[1]
					if filepath.Ext(jsURL) == ".js" {
						fullURL := url + "/" + strings.TrimPrefix(jsURL, "/")

						file.WriteString(fmt.Sprintf("%s\n", fullURL))
					}
				}
			}
		}(url)
	}

	// Wait for all goroutines to complete and print results
	wg.Wait()
	close(results)
}
