package view

import (
	"os"
	"testing"

	"coding-180/crawler/engine"
	"coding-180/crawler/frontend/model"
	common "coding-180/crawler/model"
)

func TestSearchResultView_Render(t *testing.T) {
	view := CreateSearchResultView(
		"template.html")

	out, err := os.Create("template.test.html")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: common.Profile{
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Marriage:   "离异",
			House:      "已购房",
			Hukou:      "山东菏泽",
			Education:  "大学本科",
			Cars:        "未购车",
		},
	}
	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}

	err = view.Render(out, page)
	if err != nil {
		t.Error(err)
	}

	// TODO: verify contents in template.test.html
}
