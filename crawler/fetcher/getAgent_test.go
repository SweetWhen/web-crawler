package fetcher

import (
	"testing"
	"net/http"
	"fmt"
)

func TestGetUserAgent(t * testing.T)  {
	http.HandleFunc("/", getUserAgent)
	http.ListenAndServe(":8080", nil)
}

func getUserAgent(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	fmt.Printf("user agent is: %s \n", ua)
	w.Write([]byte("user agent is " + ua))
	//Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36
}

func TestGetWithinUserAgent(t *testing.T)  {
	body, err := GET(
		"http://m.zhenai.com/zhenghun/aba",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"  )
	if err != nil {
		t.Error("get failed, %v", err)
	}
	fmt.Println(string(body))
}
