package mping

type QueuedScheduler struct {
	in         chan string
	workerChan chan chan string
}

func (q *QueuedScheduler) workChan() chan string {
	return make(chan string)
}

func (q *QueuedScheduler) submit(ip string) {
	q.in <- ip
}

func (q *QueuedScheduler) workReady(worker chan string) {
	q.workerChan <- worker
}

func (q *QueuedScheduler) run() {
	q.in = make(chan string)
	q.workerChan = make(chan chan string)
	go func() {
		var inQ []string
		var workerQ []chan string
		for {
			var activeIn string
			var activeWorker chan string
			if len(inQ) > 0 && len(workerQ) > 0 {
				activeIn = inQ[0]
				activeWorker = workerQ[0]
			}
			select {
			case in := <-q.in:
				inQ = append(inQ, in)
			case worker := <-q.workerChan:
				workerQ = append(workerQ, worker)
			case activeWorker <- activeIn:
				inQ = inQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
