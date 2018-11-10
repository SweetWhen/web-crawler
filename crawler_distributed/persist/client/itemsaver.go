package client

import (
	"log"

	"coding-180/crawler/engine"
	"coding-180/crawler_distributed/config"
	"coding-180/crawler_distributed/rpcsupport"
)

/*
初始化并连接一个rpc client到itemsaver，同时生成一个goroutine不断的从一个channel中获取要保存的item，
如果获取到了就发给itemsaver。这个函数返回这个channel
*/
func ItemSaver(
	host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++

			// Call RPC to save item
			result := ""
			err := client.Call(
				config.ItemSaverRpc,
				item, &result)

			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v: %v",
					item, err)
			}
		}
	}()

	return out, nil
}
