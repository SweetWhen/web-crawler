package parser

import (
	"regexp"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"strconv"
)

var (
	/*
	profileRe = regexp.MustCompile(
		`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(
		`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
	*/
	profileRe  = regexp.MustCompile(`<a href="(http://m.zhenai.com/u/[0-9]+)#seo" class="left-item" data-v-d9d4d86c>`)
	cityUrlRe = regexp.MustCompile(` href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
)
var peopleCnt int
func ParseCity(
	contents []byte, _ string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(
		contents, -1)
	//fmt.Println("********************************************")
	//fmt.Println(string(contents))
	//fmt.Println("********************************************")
	result := engine.ParseResult{}
	for _, m := range matches {
		peopleCnt++
		result.Requests = append(
			result.Requests, engine.Request{
				Url: string(m[1]),
				Parser: NewProfileParser(
					strconv.Itoa(peopleCnt)),
			})
	}

	matches = cityUrlRe.FindAllSubmatch(
		contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests,
			engine.Request{
				Url: string(m[1]),
				Parser: engine.NewFuncParser(
					ParseCity, config.ParseCity),
			})
	}

	return result
}
