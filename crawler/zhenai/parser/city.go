package parser

import (
	"regexp"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"fmt"
)

var (

	profileRe = regexp.MustCompile(
		`<a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank">[^<]+</a></th></tr><tr><td width="180"><span class="grayL">性别：</span>([^<]+)</td>`)
	cityUrlRe = regexp.MustCompile(
		`href="(http://www.zhenai.com/zhenghun/[a-z]+/[0-9]+)"`)
	/*
	profileRe  = regexp.MustCompile(`<a href="(http://m.zhenai.com/u/[0-9]+)#seo" class="left-item" data-v-d9d4d86c>`)
	cityUrlRe = regexp.MustCompile(` href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
	*/

//<a href="http://album.zhenai.com/u/1486691188" target="_blank">回忆很伤感</a></th></tr><tr><td width="180"><span class="grayL">性别：</span>男士</td>
//a href="http://www.zhenai.com/zhenghun/foshan/3">3</a> <!----></li>
)

func ParseCity(
	contents []byte, _ string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(
		contents, -1)
fmt.Println(string(contents))
	result := engine.ParseResult{}
	for _, m := range matches {
		fmt.Println("********************************************")
		fmt.Println(string(m[0]))
		fmt.Println(string(m[1]))
		fmt.Println(string(m[2]))


		result.Requests = append(
			result.Requests, engine.Request{
				Url: string(m[1]),
				Parser: NewProfileParser(
					string(m[2])),
			})
	}

	matches = cityUrlRe.FindAllSubmatch(
		contents, -1)
	for _, m := range matches {
		fmt.Println("next city:", string(m[1]))
		fmt.Println("********************************************")
		result.Requests = append(result.Requests,
			engine.Request{
				Url: string(m[1]),
				Parser: engine.NewFuncParser(
					ParseCity, config.ParseCity),
			})
	}

	return result
}
