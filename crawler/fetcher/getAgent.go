package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bufio"
	"golang.org/x/text/transform"
)

func GET(url, userAgent string ) ([]byte, error){
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("User-Agent", userAgent)

	resp,err := client.Do(request)
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
