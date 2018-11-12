package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"

	"log"

	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"coding-180/crawler/config"
)

var (
	rateLimiter = time.Tick(
		time.Second / config.Qps)
	verboseLogging = false
)

func SetVerboseLogging() {
	verboseLogging = true
}

func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	if verboseLogging {
		log.Printf("Fetching url %s", url)
	}
	//resp, err := http.Get(url)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println("======================================================")
	fmt.Println(url)
	request.Header.Set(
		"Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set(
		"Upgrade-Insecure-Requests",
		"1")
	request.Header.Set(
		"Proxy-Connection",
		"keep-alive")
	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36" )

	resp,err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf("wrong status code: %d",
				resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader,
		e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(
	r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(
		bytes, "")
	return e
}
