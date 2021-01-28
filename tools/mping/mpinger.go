package mping

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"log"
	"mytools/tools/ipz"
	"mytools/utils"
	"runtime"
	"time"
)

type Pingers struct {
	workerCount int
	repeat      int
	spacing     int
	IPTrunk     IPTrunk
	scheduler   Scheduler
	pinger      pinger
}

type Scheduler interface {
	workChan() chan string
	run()
	submit(string)
	workReady(chan string)
}

type pinger struct {
	count     int
	interval  time.Duration
	privilege bool
}

type IPTrunk interface {
	Output(chan string)
	Len() uint32
}

func NewPingers(count int, spacing int, workerCount int, repeat int,
	interval int, ipfile string, ipnet string) Pingers {
	var iplist IPTrunk
	var privilege bool
	if ipfile != "" {
		iplist = newFileIP(ipfile)
	} else if ipnet != "" {
		iplist = ipz.NewIPRange(ipnet)
	} else {
		panic("请使用ipnet或ipfile提供测试IP地址！")
	}

	if sysType := runtime.GOOS; sysType == "windows" {
		privilege = true
	}
	return Pingers{
		workerCount: workerCount,
		repeat:      repeat,
		spacing:     spacing,
		scheduler:   &QueuedScheduler{},
		IPTrunk:     iplist,
		pinger: pinger{
			count,
			time.Duration(interval),
			privilege,
		},
	}
}

func (p *Pingers) Run() {
	ipNumber := p.IPTrunk.Len()
	if ipNumber == 0 {
		fmt.Println("没有可执行的IP列表！")
		return
	}
	p.scheduler.run()
	out := make(chan *ping.Statistics)
	for i := 0; i < p.workerCount; i++ {
		createWorker(p.scheduler.workChan(), out, p.pinger, p.scheduler)
	}
	for i := 0; i < p.repeat; i++ {
		x := i + 1
		go func() {
			out := make(chan string)
			p.IPTrunk.Output(out)
			for {
				ip := <-out
				p.scheduler.submit(ip)
			}
		}()
		fmt.Printf("==================第%v次测试=================\n", x)
		var j uint32 = 0
		for ; j < ipNumber; j++ {
			stat := <-out
			message := fmt.Sprintf("%v ping statistics: %v packets transmitted, %v received, %v%% packet loss, max %v\n",
				stat.Addr, stat.PacketsSent, stat.PacketsRecv, stat.PacketLoss, stat.MaxRtt)
			if stat.PacketLoss != 0 {
				c := utils.NewColor(1, 40, 31)
				c.Run(message)
			} else {
				c := utils.NewColor(1, 40, 32)
				c.Run(message)
			}
		}
		if x != p.repeat {
			time.Sleep(time.Duration(p.spacing) * time.Second)
		}
	}
}

func createWorker(in chan string, out chan *ping.Statistics, pinger pinger, s Scheduler) {
	go func() {
		for {
			s.workReady(in)
			ip := <-in
			stats, err := pingWorker(ip, pinger.count, pinger.interval, pinger.privilege)
			if err != nil {
				log.Println(err)
				continue
			}
			out <- stats
		}
	}()
}

func pingWorker(ip string, count int, interval time.Duration, privilege bool) (stats *ping.Statistics, err error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		//message := fmt.Sprintf("Make pinger of %v is error!\n", ip)
		//return stats, errors.New(message)
		return stats, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = count
	pinger.Interval = interval * time.Millisecond
	pinger.Timeout = time.Duration(3*count) * pinger.Interval
	pinger.Run()
	stats = pinger.Statistics()
	return
}
