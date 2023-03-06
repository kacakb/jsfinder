package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {

	var urlsFilePath string
	flag.StringVar(&urlsFilePath, "l", "", "filename to read URLS from")

	var limit int = 20
	flag.IntVar(&limit, "c", limit, "concurrency limit")

	var silent bool
	flag.BoolVar(&silent, "s", false, "Program is running in silent mode")

	var outputFile string
	flag.StringVar(&outputFile, "o", "output.txt", "Filename to write found URLs to")

	var readURLs bool
	flag.BoolVar(&readURLs, "read", false, "Read URLs from stdin")

	flag.Parse()

	if readURLs && urlsFilePath != "" {
		fmt.Println("Cannot use -l and -read flags together")
		return
	}

	if readURLs {
		urlsFilePath = os.Stdin.Name()
	}

	if !readURLs && urlsFilePath == "" {
		fmt.Println("Please provide a file containing URLs with the -l flag or use -read flag to read from stdin")
		return

	} else if readURLs && urlsFilePath == "" {
		fmt.Println("Please provide a list of URLs to scan:")
		return
	} else {
		fmt.Println("Reading URLs from file:", urlsFilePath)
	}

	if limit == 20 {
		fmt.Println("Concurrency limit is running default: 20")
		if !silent {
			fmt.Println("Verbose mode active")
		} else {
			fmt.Println("Silent mode active")
		}
	} else {
		if !silent {
			fmt.Println("Verbose mode active")
			fmt.Printf("Concurrency limit is running %d\n", limit)
		} else {
			fmt.Printf("Silent mode active\nConcurrency limit is running %d\n", limit)
		}
	}

	var urlsFile *os.File
	var err error

	if readURLs {
		urlsFile = os.Stdin
	} else {
		urlsFile, err = os.Open(urlsFilePath)
		if err != nil {
			panic(err)
		}
		defer urlsFile.Close()
	}

	results := make(chan string)

	sem := make(chan struct{}, limit)

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(urlsFile)
	for scanner.Scan() {
		url := scanner.Text()

		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: 7 * time.Second,
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				if !silent {
					fmt.Printf("Error creating request for %s: %v\n", url, err)
				}
				return

			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")

			resp, err := client.Do(req)
			if err != nil {
				if !silent {
					fmt.Printf("Error getting response from %s: %v\n", url, err)
				}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				if err != nil {
					if !silent {
						fmt.Printf("Error getting response from %s: status code %d\n", url, resp.StatusCode)
					}
				}
				return
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				if !silent {
					fmt.Println("Error getting response from", url, "status code", resp.StatusCode)
				}
				return
			}
			bodyString := string(bodyBytes)

			re := regexp.MustCompile(`(?i)(?:src|srcdoc|formaction|dynsrc|standby|ng-include|ui-sref|href|data-main|data|onclick|onload|style|srcdoc|formaction|iframe|object|background|input|button|action|dynsrc|srcset|manifest|code|archive|classid|cite|codebase|longdesc|lowsrc|usemap|standby|ng-click|ng-src|ng-inlude|ui-sref|require|createElement|appendChild|innerHTML|getScript|XMLHttpRequest|fetch|import|onerror|WebSocket|ServiceWorker|SharedWorker|importScripts|eval)\s*=\s*["']([^"']*\.js(\?[^"']*)?)["']`)

			matches := re.FindAllStringSubmatch(bodyString, -1)
			if len(matches) > 0 {
				var file *os.File
				var err error

				if outputFile != "" {
					file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Printf("Error opening file %s for writing: %v\n", outputFile, err)
						return
					}
				} else {
					file, err = os.OpenFile("found.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Printf("Error opening file for writing: %v\n", err)
						return
					}
				}
				defer file.Close()

				for _, match := range matches {
					jsURL := match[1]
					if strings.HasSuffix(jsURL, ".js") || strings.Contains(jsURL, ".js?") {

						if strings.HasPrefix(jsURL, "/") {
							if strings.Contains(url, ".com") || strings.Contains(url, ".net") || strings.Contains(url, ".org") {
								if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
									url = "https://" + strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")
								}
								if strings.Contains(jsURL, ".com") || strings.Contains(jsURL, ".net") || strings.Contains(jsURL, ".org") {
									if strings.HasPrefix(jsURL, "//") {
										file.WriteString(fmt.Sprintf("https:%s\n", jsURL))
									} else {
										file.WriteString(fmt.Sprintf("https://%s\n", strings.TrimPrefix(jsURL, "/")))
									}
								} else {
									if strings.HasPrefix(jsURL, "/") {
										file.WriteString(fmt.Sprintf("%s%s\n", url, jsURL))
									} else if strings.HasPrefix(jsURL, "https://") || strings.HasPrefix(jsURL, "http://") {
										file.WriteString(fmt.Sprintf("%s\n", jsURL))
									} else if strings.HasPrefix(jsURL, "//") {
										file.WriteString(fmt.Sprintf("https:%s\n", jsURL))
									} else {
										file.WriteString(fmt.Sprintf("https://%s\n", jsURL))
									}
								}
							} else {
								file.WriteString(fmt.Sprintf("%s/%s\n", url, jsURL))
							}
						} else if strings.HasPrefix(jsURL, "https://") || strings.HasPrefix(jsURL, "http://") {
							file.WriteString(fmt.Sprintf("%s\n", jsURL))
						} else if strings.Contains(jsURL, ".com") || strings.Contains(jsURL, ".net") || strings.Contains(jsURL, ".org") {
							if strings.Contains(jsURL, ".com/") {
								file.WriteString(fmt.Sprintf("https://%s%s\n", jsURL[:strings.Index(jsURL, ".com")+4], jsURL[strings.Index(jsURL, ".com")+4:]))
							} else if strings.Contains(jsURL, ".net/") {
								file.WriteString(fmt.Sprintf("https://%s%s\n", jsURL[:strings.Index(jsURL, ".net")+4], jsURL[strings.Index(jsURL, ".net")+4:]))
							} else if strings.Contains(jsURL, ".org/") {
								file.WriteString(fmt.Sprintf("https://%s%s\n", jsURL[:strings.Index(jsURL, ".org")+4], jsURL[strings.Index(jsURL, ".org")+4:]))
							} else {
								file.WriteString(fmt.Sprintf("https://%s/%s\n", jsURL[:strings.Index(jsURL, ".")+4], jsURL[strings.Index(jsURL, ".")+4:]))
							}
						} else {
							file.WriteString(fmt.Sprintf("%s/%s\n", url, jsURL))
						}
					} else {
						file.WriteString(fmt.Sprintf("%s/%s\n", url, jsURL))
					}
				}

			}
		}(url)
	}

	wg.Wait()
	close(results)
}
