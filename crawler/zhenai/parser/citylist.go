package parser

import (
	"regexp"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
)

//const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`
const cityListRe = `{linkContent:"([^"]+)",linkURL:"(http://m.zhenai.com/zhenghun/[0-9a-z]+)"}`
func ParseCityList(
	contents []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, m := range matches {
		result.Requests = append(
			result.Requests, engine.Request{
				Url: string(m[2]),
				Parser: engine.NewFuncParser(
					ParseCity, config.ParseCity),
			})
	}

	return result
}
