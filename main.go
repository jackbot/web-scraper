package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func readUrls(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getTitleFromUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	body, err := getBody(url)

	if err != nil {
		fmt.Printf("Couldn't connect to site: %s\n", url)
		return
	}

	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		fmt.Printf("Couldn't parse HTML: %s\n", url)
		return
	}

	doc.Find("html").Each(func(i int, s *goquery.Selection) {
		title := s.Find("title").Text()
		fmt.Printf("%s -> %s\n", url, title)
	})
}

func getBody(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 status code: %d", res.StatusCode)
	}

	return res.Body, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run main.go /path/to/url_file\n")
		return
	}

	path := os.Args[1]

	urls, err := readUrls(path)

	if err != nil {
		fmt.Printf("Couldn't read URLs from file: %s\n", path)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go getTitleFromUrl(url, &wg)
	}

	wg.Wait()
}
