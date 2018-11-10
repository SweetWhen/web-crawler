package worker

import "coding-180/crawler/engine"

type CrawlService struct{}

func (CrawlService) Process(
	req Request, result *ParseResult) error {
	engineReq, err := DeserializeRequest(req)
	if err != nil {
		return err
	}

	engineResult, err := engine.Worker(engineReq)
	if err != nil {
		return err
	}

	//简单的给result赋值就把结果给返回给client了...
	*result = SerializeResult(engineResult)
	return nil
}
