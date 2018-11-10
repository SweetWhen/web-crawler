package main

import (
	"coding-180/crawler/config"
	"coding-180/crawler/engine"
	"coding-180/crawler/persist"
	"coding-180/crawler/scheduler"
	"coding-180/crawler/zhenai/parser"
	"flag"
	"log"
)
var dbAddr = flag.String("db_addr", "http://127.0.0.1:9200", "db addr, the db must up")

func main() {
	flag.Parse()
	itemChan, err := persist.ItemSaver(*dbAddr,
		config.ElasticIndex)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to db addr: %s\n", *dbAddr )
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList),
	})
}
