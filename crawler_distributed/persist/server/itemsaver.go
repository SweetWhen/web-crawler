package main

import (
	"log"

	"fmt"

	"flag"

	"gopkg.in/olivere/elastic.v5"
	"coding-180/crawler/config"
	"coding-180/crawler_distributed/persist"
	"coding-180/crawler_distributed/rpcsupport"
	"strings"
)

var port = flag.Int("port", 0,
	"the port for me to listen on")
var dbAddr = flag.String("db_addr",  "http://127.0.0.1:9200",
				"saver ElasticSearch addr")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	if 0 != strings.Index(*dbAddr, "http://") {
		fmt.Println("db_addr format http://[ip addr]:[port number]")
		return
	}
	log.Fatal(serveRpc(
		fmt.Sprintf(":%d", *port),
		*dbAddr,
		config.ElasticIndex))
}

func serveRpc(host, dbAddr, index string) error {
	client, err := elastic.NewClient(
		elastic.SetURL(dbAddr),
		elastic.SetSniff(false))
	if err != nil {
		return err
	}

	return rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client: client,
			Index:  index,
		})
}
