package main

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func main() {
	numThreads := 10000
	file, err := os.Open("wp-sites.txt")

	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()

	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	c := colly.NewCollector(colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: numThreads})
	defaultRoundTripper := http.DefaultTransport
	defaultRoundTripperPointer := defaultRoundTripper.(*http.Transport)
	t := *defaultRoundTripperPointer // deref to get copy
	t.MaxIdleConns = 0
	t.MaxIdleConnsPerHost = 1
	c.WithTransport(&t)
	smily := html.UnescapeString("&#" + strconv.Itoa(0x1F60A) + ";")
	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "wp") {
			if _, err := fmt.Fprintln(w, r.Request.URL.Host+smily); err != nil {
				panic(err)
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	q, _ := queue.New(
		numThreads, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: len(text)}, // Use default queue storage
	)
	for _, each_host := range text {
		q.AddURL("http://" + each_host)
	}
	q.Run(c)
	c.Wait()
	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}

}
