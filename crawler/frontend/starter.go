package main

import (
	"net/http"

	"coding-180/crawler/frontend/controller"
	"fmt"
)
const webRoot = "/opt/mygo/src/coding-180/crawler/frontend/view"
func main() {
	//设置默认页面路径为crawler/frontend/view
	http.Handle("/", http.FileServer(http.Dir(webRoot)))
	//search路由
	http.Handle("/search",
		controller.CreateSearchResultHandler(
			fmt.Sprintf("%s%s", webRoot, "/template.html")))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
