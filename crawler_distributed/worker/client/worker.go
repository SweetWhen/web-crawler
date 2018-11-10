package client

import (
	"net/rpc"

	"coding-180/crawler/engine"
	"coding-180/crawler_distributed/config"
	"coding-180/crawler_distributed/worker"
)

/*
返回一个函数，这个函数会从 worker rpc clients channel中获取一个rpc client，然后将一个req序列化成网络可以传播的
字符串形式，并通过rpc发送给被选中的worker，由worker来处理这个req。
*/
func CreateProcessor(
	clientChan chan *rpc.Client) engine.Processor {

	return func(
		req engine.Request) (
		engine.ParseResult, error) {

		sReq := worker.SerializeRequest(req)

		var sResult worker.ParseResult
		c := <-clientChan
		err := c.Call(config.CrawlServiceRpc,
			sReq, &sResult)

		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult),
			nil
	}
}
