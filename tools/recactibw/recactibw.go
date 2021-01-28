package recactibw

import (
	"fmt"
	"mytools/utils"
	"time"
)

type Engine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Repair struct {
	StartTime int64
	EndTime   int64
	Interval  int64
	Percent   float64
	Filepath  string
	LogChan   chan string
	Details   bool
}

type Scheduler interface {
	Run()
	WorkChan() chan Repair   //生成一个in chan，没有其实也行
	WorkerReady(chan Repair) //接受worker in chan的chan
	Submit(Repair)
}

func (e *Engine) Run(startTime, endTime, interval int64,
	percent float64, source string, details bool, logfile string, print bool) {
	// 获取所有需要改动的文件

	utils.ClearScreen()
	time.Sleep(1 * time.Second)

	out := make(chan string)
	e.Scheduler.Run()

	myBarChans := []utils.MyBarChan{}
	for i := 0; i < e.WorkerCount+1; i++ {
		barChan := utils.MyBarChan{
			AddChan: make(chan int),
			OptChan: make(chan utils.MyBarOptChan),
		}
		utils.NewMyProgressBars(barChan, i)
		myBarChans = append(myBarChans, barChan)
	}

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkChan(), out, print, myBarChans[i+1])
	}

	errLogIn := make(chan string)
	utils.MyLoger("/var/log/mytools/errors.log", errLogIn, false)

	files := utils.GetAllFiles(source)
	for _, f := range files {
		r := Repair{
			StartTime: startTime,
			EndTime:   endTime,
			Interval:  interval,
			Percent:   percent,
			Filepath:  f,
			LogChan:   errLogIn,
			Details:   details,
		}
		e.Scheduler.Submit(r)
	}

	logInChan := make(chan string)
	utils.MyLoger(logfile, logInChan, details)

	myBarChans[0].OptChan <- utils.MyBarOptChan{
		Desc:  "处理文件总进度",
		Total: len(files),
	}
	for i := 0; i < len(files); i++ {
		myBarChans[0].AddChan <- 1
		log := <-out
		logInChan <- log
	}
}

func (e *Engine) createWorker(in chan Repair, out chan string, print bool, bar utils.MyBarChan) (err error) {
	go func() {
		for {
			e.Scheduler.WorkerReady(in)
			repair := <-in
			if print {
				err = RepairWorkerFunc(repair, bar, Print)
			} else {
				err = RepairWorkerFunc(repair, bar, Modify)
			}
			if err != nil {
				repair.LogChan <- err.Error()
			}
			res := fmt.Sprintf("完成操作文件: %v\n", repair.Filepath)
			out <- res
		}
	}()
	return nil
}
