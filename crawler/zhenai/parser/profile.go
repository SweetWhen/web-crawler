package parser

import (
	"regexp"
	"strconv"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"coding-180/crawler/model"

)

var userName = regexp.MustCompile(`<span class="nick_name">昵称</span>\s*<span>([^<]+)</span>` )
var educationRe = regexp.MustCompile(`<span class="nick_name">学历</span>\s*<span>([^<]+)</span>` )
var ageRe = regexp.MustCompile(`<span class="nick_name">年龄</span>\s*<span>([\d]+)岁</span>`)
var heightRe = regexp.MustCompile(`<span class="nick_name">身高</span>\s*<span>([\d]+)CM</span>`)
var weightRe = regexp.MustCompile(`<span class="nick_name">体重</span>\s*<span>([\d]+)KG</span>`)
var incomeRe = regexp.MustCompile(`<span class="nick_name">月收入</span>\s*<span>([^<]+)</span>`)
var marriageRe = regexp.MustCompile(`<span class="nick_name">婚姻状况</span>\s*<span>([^<]+)</span>`)
var hukouRe = regexp.MustCompile(`<span class="nick_name">籍贯</span>\s*<span>([^<]+)</span>`)
var houseRe = regexp.MustCompile(`<span class="nick_name">住房情况</span>\s*<span>([^<]+)</span>`)
var carsRe = regexp.MustCompile(`<span class="nick_name">买车情况</span>\s*<span>([^<]+)</span>`)
var urlIdRe =  regexp.MustCompile(`http://m.zhenai.com/u/([\d]+)`)

func parseProfile(
	contents []byte, url string,
	_ string) engine.ParseResult {
	profile := model.Profile{}
	profile.Name = extractString(contents,userName)
	profile.Education = extractString(contents,educationRe)
	//fmt.Println("********************************************")
	//fmt.Println(string(contents))
	//fmt.Println("********************************************")


	age, err :=  strconv.Atoi(
		extractString(contents,ageRe) )
	if err == nil {
		profile.Age = age
	}
	height, err :=  strconv.Atoi(
		extractString(contents,heightRe) )
	if err == nil {
		profile.Height = height
	}
	weight, err :=  strconv.Atoi(
		extractString(contents,weightRe) )
	if err == nil {
		profile.Weight = weight
	}
	profile.Income = extractString(contents,incomeRe)
	profile.Marriage = extractString(contents,marriageRe)
	profile.Hukou = extractString(contents,hukouRe)
	profile.House = extractString(contents,houseRe)
	profile.Cars = extractString(contents,carsRe)

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:  url,
				Type: "zhenai",
				Id: extractString(
					[]byte(url), urlIdRe),
				Payload: profile,
			},
		},
	}

	return result
}

func extractString(
	contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(
	contents []byte,
	url string) engine.ParseResult {
	return parseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (
	name string, args interface{}) {
	return config.ParseProfile, p.userName
}

func NewProfileParser(
	name string) *ProfileParser {
	return &ProfileParser{
		userName: name,
	}
}
