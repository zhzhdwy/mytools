package recactibw

type QueueScheduler struct {
	RepairChan chan Repair
	WorkerChan chan chan Repair
}

func (q *QueueScheduler) WorkChan() chan Repair {
	return make(chan Repair)
}

func (q *QueueScheduler) Submit(repair Repair) {
	q.RepairChan <- repair
}

func (q *QueueScheduler) WorkerReady(in chan Repair) {
	q.WorkerChan <- in
}

func (q *QueueScheduler) Run() {
	q.RepairChan = make(chan Repair)
	q.WorkerChan = make(chan chan Repair)

	go func() {
		//t := time.Tick(30 * time.Second)
		var reqairQ []Repair
		var workerQ []chan Repair
		for {
			var activeReqair Repair      //这个代表外部传入
			var activeWorker chan Repair //这个代表worker的接受通道
			if len(reqairQ) > 0 && len(workerQ) > 0 {
				activeReqair = reqairQ[0]
				activeWorker = workerQ[0]
			}
			select {
			case w := <-q.WorkerChan:
				workerQ = append(workerQ, w)
			case r := <-q.RepairChan:
				reqairQ = append(reqairQ, r)
			case activeWorker <- activeReqair:
				reqairQ = reqairQ[1:]
				workerQ = workerQ[1:]
				//case <-t:
				//	log.Printf("reqairQ: %v, workerQ: %v\n", len(reqairQ), len(workerQ))
			}
		}
	}()
}
