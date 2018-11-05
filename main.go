// test_go_mail.ru project main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"test_go_mail.ru/counter"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("Hello,separate URLs by comma please")
	wordToFind := flag.String("q", "go", "word to find from url")
	k := flag.Int("k", 5, "number of gouroutines")
	flag.Parse()

	totalURLsCount()
	parser(*k, wordToFind, findStartOfUrls())

}

func checkURL(url string) bool {
	_, err := http.Get(url)
	if err != nil {
		return false
	} else {
		return true
	}
}

func totalURLsCount() {
	urlArgs := ""

	for _, args := range os.Args {
		if strings.HasPrefix(args, "http") == true {
			urlArgs = args
		}
	}
	totalURL := 0
	for _, url := range strings.Split(urlArgs, ",") {
		if checkURL(url) == true {
			totalURL++
		}
	}
	fmt.Println("Total URLs : ", totalURL)
}

func parser(k int, wordToFind *string, urlArgs string) {

	var wg sync.WaitGroup
	ch := make(chan struct{}, k)
	tc := counter.NewCounter()

	for _, url := range strings.Split(urlArgs, ",") {
		if checkURL(url) == true {

			wg.Add(1)
			ch <- struct{}{}

			go parseSite(url, wordToFind, &wg, tc)
			<-ch

		} else {
			fmt.Println(url, "--- cant access this URL")
		}
	}
	wg.Wait()
	fmt.Println("Total:", tc.Total)
	fmt.Println("Done")
}

func findStartOfUrls() string {
	urlArgs := ""

	for _, args := range os.Args {
		if strings.HasPrefix(args, "http") == true {
			urlArgs = args
		}
	}
	return urlArgs
}

func parseSite(url string, wordToFind *string, wg *sync.WaitGroup, tc *counter.TotalCounter) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(url, " - cant access this URL")

	}
	bytes, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	counter := 0
	text := strings.Fields(string(bytes))
	for _, word := range text {
		if strings.Contains(strings.ToLower(word), *wordToFind) == true {
			counter++
		}

	}
	tc.SafeAdd(counter)
	fmt.Println("Counts for: ", url, " -- ", counter)

}
