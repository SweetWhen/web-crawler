package scheduler

import "coding-180/crawler/engine"

type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request
}
//生成一个channel
func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}
//这个方法只是把request赋值给requestChan，Schedule模块会把这些request保存到一个队列中去。
func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}
//如果一个worker闲下来了调用这个方法把它送到workerChan中，Schedule模块会去搜集这些闲下来的worker
func (s *QueuedScheduler) WorkerReady(
	w chan engine.Request) {
	s.workerChan <- w
}

//生成workerChan和requestChan，并用一个goroutine不断的从这两个channel收东西加到队列中，
//并各取出一个req和worker，将这个req发给这个worker。
func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 &&
				len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}

			select {
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
