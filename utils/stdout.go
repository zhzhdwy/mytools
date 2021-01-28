package utils

import (
	"fmt"
	"log"
	"runtime"
)

type color struct {
	d int
	b int
	f int
}

func NewColor(d int, b int, f int) color {
	return color{
		d: d,
		b: b,
		f: f,
	}
}

func (c *color) Run(message string) {
	if sysType := runtime.GOOS; sysType == "linux" {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, c.d, c.b, c.f, message, 0x1B)
		return
	}
	log.Printf("%v", message)
}
