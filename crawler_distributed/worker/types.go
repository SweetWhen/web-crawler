package worker

import (
	"errors"

	"fmt"
	"log"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"coding-180/crawler/zhenai/parser"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

func SerializeResult(
	r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		result.Requests = append(result.Requests,
			SerializeRequest(req))
	}
	return result
}

func DeserializeRequest(
	r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

func DeserializeResult(
	r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error deserializing "+
				"request: %v", err)
			continue
		}
		result.Requests = append(result.Requests,
			engineReq)
	}
	return result
}

//返回值中的engine.Parser是一个interface，这个interface只要某个对象实现了它的方法它就可以接受
//对象里面有一个函数属性，这里直接把一个函数当成一等公民（类似一个变量）赋给一个对象的属性了，非常之灵活。
func deserializeParser(
	p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		return engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList), nil
	case config.ParseCity:
		return engine.NewFuncParser(
			parser.ParseCity,
			config.ParseCity), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	case config.ParseProfile:
		if userName, ok := p.Args.(string); ok {
			return parser.NewProfileParser(
				userName), nil
		} else {
			return nil, fmt.Errorf("invalid "+
				"arg: %v", p.Args)
		}
	default:
		return nil, errors.New(
			"unknown parser name")
	}
}
