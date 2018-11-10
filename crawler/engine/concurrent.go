package engine

type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	//让schedule跑起来
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(),
			out, e.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			continue
		}
		e.Scheduler.Submit(r)
	}

	for {
		/*
		所有的worker产生的结果都通过这个out返回，这里engine不断的从这个out中读取结果。
		*/
		result := <-out
		for _, item := range result.Items {
			//item是我们最终想要保存的有用信息，起一个goroutine把它写到数据库
			go func(i Item) {
				e.ItemChan <- i
			}(item)
		}

		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)//把这个request提交到schedule，让它来安排由哪一个worker处理。
		}
	}
}
/*
生成一个goroutine，不断的从in 读取req，然后调用e.RequestProcessor处理req，并等待其返回，
返回之后把结果发送到out，然后又把这个goroutine提交到schedule，并用 <-in阻塞等待其被调用
*/
func (e *ConcurrentEngine) createWorker(
	in chan Request,
	out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in

			/*
			e.RequestProcessor这个会从 worker client channel获取一个client，然后将req通过rpc发送给它，并等待它返回
			这里等待不要紧，因为这是在一个goroutine中阻塞等待的。
			 */
			result, err := e.RequestProcessor(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}
