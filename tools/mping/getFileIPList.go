package mping

import (
	"io/ioutil"
	"strings"
)

type fileIP struct {
	iplist []string
}

func newFileIP(filePath string) *fileIP {
	ips, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	list := strings.Split(string(ips), "\n")
	var iplist []string
	for _, ip := range list {
		ip := strings.Replace(ip, "\r", "", -1)
		if ip != "" {
			iplist = append(iplist, ip)
		}
	}
	return &fileIP{iplist: iplist}
}

func (f fileIP) Output(out chan string) {
	go func() {
		for _, ip := range f.iplist {
			out <- ip
		}
	}()
}

func (f fileIP) Len() uint32 {
	listLen := uint32(len(f.iplist))
	return listLen
}
