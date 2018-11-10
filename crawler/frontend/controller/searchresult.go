package controller

import (
	"net/http"
	"strconv"
	"strings"

	"context"
	"reflect"

	"regexp"

	"gopkg.in/olivere/elastic.v5"
	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"coding-180/crawler/frontend/model"
	"coding-180/crawler/frontend/view"
)

type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client
}

func CreateSearchResultHandler(
	template string) SearchResultHandler {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return SearchResultHandler{
		view: view.CreateSearchResultView(
			template),
		client: client,
	}
}

//client的请求中带了/search的话会到这个函数中。
func (h SearchResultHandler) ServeHTTP(
	w http.ResponseWriter, req *http.Request) {
		//从request中拿q参数
	q := strings.TrimSpace(req.FormValue("q"))
	//从request中拿from参数
	from, err := strconv.Atoi(
		req.FormValue("from"))
	if err != nil {
		from = 0
	}
//根据url去数据库拿记录
	page, err := h.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	err = h.view.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
		return
	}
}

const pageSize = 10

func (h SearchResultHandler) getSearchResult(
	q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q
//向数据库请求数据。。。
	tmpQ := q
	if strings.Contains(tmpQ, ":"){
		tmpQ = rewriteQueryString(tmpQ)
	}
	if tmpQ == "" {
		tmpQ = "*"
	}
	resp, err := h.client.
		Search(config.ElasticIndex).
		Query(elastic.NewQueryStringQuery(tmpQ)).
		From(from).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(
		reflect.TypeOf(engine.Item{}))
	if result.Start == 0 {
		result.PrevFrom = -1
	} else {
		result.PrevFrom =
			(result.Start - 1) /
				pageSize * pageSize
	}
	result.NextFrom =
		result.Start + len(result.Items)

	return result, nil
}

// Rewrites query string. Replaces field names
// like "Age" to "Payload.Age"
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}
