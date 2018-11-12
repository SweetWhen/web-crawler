package main

import (
	"errors"
	"net/rpc"

	"log"

	"flag"

	"strings"

	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"coding-180/crawler/scheduler"
	"coding-180/crawler/zhenai/parser"
	itemsaver "coding-180/crawler_distributed/persist/client"
	"coding-180/crawler_distributed/rpcsupport"
	worker "coding-180/crawler_distributed/worker/client"
	"fmt"
)

var (
	itemSaverHost = flag.String(
		"itemsaver_host", "", "itemsaver host: 127.0.0.1:1234")

	workerHosts = flag.String(
		"worker_hosts", "",
		"worker hosts (comma separated): 127.0.0.1:9000,127.0.0.1:9001")
)

func main() {
	flag.Parse()
	//得到保存item的channel
	if *itemSaverHost == "" {
		fmt.Println("must specify a itemsaver host!")
		return
	}
	if *workerHosts == "" {
		fmt.Println("must specify a worker host!")
		return
	}
	itemChan, err := itemsaver.ItemSaver(
		*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool, err := createClientPool(
		strings.Split(*workerHosts, ","))
	if err != nil {
		panic(err)
	}
	//得到req处理函数
	processor := worker.CreateProcessor(pool)

	e := engine.ConcurrentEngine{
		//Scheduler是一个interface，只要实现了它方法的对象，它都可以接受。
		Scheduler:        &scheduler.QueuedScheduler{}, //这里使用了QueueScheduler这个类型的对象，注意它的几个方法。
		WorkerCount:      100,
		ItemChan:         itemChan, //保存item的channel放在这里
		RequestProcessor: processor, //req处理函数保存在这里
	}

	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		//Parser是一个interface，需要实现了它方法的对象。
		Parser: engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList),
	})
}

/*
这个函数根据cli给的worker host IP列表生成多个rpc client，并将这些rpc clients装到一个队列clients中，
同时生成一个goroutine循环不断的将这些client塞进一个channel，这个channel会返回给engine，engine就是从这个channel
不断的拿rpc client，并把request发给这些worker。
*/
func createClientPool(
	hosts []string) (chan *rpc.Client, error) {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Connected to %s", h)
		} else {
			log.Printf(
				"Error connecting to %s: %v",
				h, err)
		}
	}

	if len(clients) == 0 {
		return nil, errors.New(
			"no connections available")
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out, nil
}
