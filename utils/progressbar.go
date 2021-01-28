package utils

import (
	"github.com/schollz/progressbar/v3"
	"sync"
)

type MyBarChan struct {
	AddChan chan int
	OptChan chan MyBarOptChan
}

type MyBarOptChan struct {
	Desc  string
	Total int
}

var m sync.Mutex

//line是从屏幕下方网上算
func NewMyProgressBars(barChan MyBarChan, line int) {
	bar := progressbar.NewOptions(1000,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	go func() {
		for {
			select {
			case a := <-barChan.AddChan:
				m.Lock()
				CursorDown(line)
				bar.Add(a)
				CursorUp(line)
				m.Unlock()
			case desc := <-barChan.OptChan:
				bar.Reset()
				bar.Describe(desc.Desc)
				bar.ChangeMax(desc.Total)
			}
		}
	}()
}
