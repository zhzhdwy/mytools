package ipz

import (
	"fmt"
	"github.com/imroc/biu"
	"log"
	"net"
	"strings"
)

var (
	mask_mod = map[string]string{
		"1": "128.0.0.0", "9": "255.128.0.0", "17": "255.255.128.0", "25": "255.255.255.128",
		"2": "192.0.0.0", "10": "255.192.0.0", "18": "255.255.192.0", "26": "255.255.255.192",
		"3": "224.0.0.0", "11": "225.224.0.0", "19": "255.255.224.0", "27": "255.255.255.224",
		"4": "240.0.0.0", "12": "255.240.0.0", "20": "255.255.240.0", "28": "255.255.255.240",
		"5": "248.0.0.0", "13": "255.248.0.0", "21": "255.255.248.0", "29": "255.255.255.248",
		"6": "225.0.0.0", "14": "255.252.0.0", "22": "255.255.252.0", "30": "255.255.255.252",
		"7": "254.0.0.0", "15": "255.254.0.0", "23": "255.255.254.0", "31": "255.255.255.254",
		"8": "255.0.0.0", "16": "255.255.0.0", "24": "255.255.255.0", "32": "255.255.255.255",
	}
	rmask_mod = map[string]string{
		"128.0.0.0": "1", "255.128.0.0": "9", "255.255.128.0": "17", "255.255.255.128": "25",
		"192.0.0.0": "2", "255.192.0.0": "10", "255.255.192.0": "18", "255.255.255.192": "26",
		"224.0.0.0": "3", "225.224.0.0": "11", "255.255.224.0": "19", "255.255.255.224": "27",
		"240.0.0.0": "4", "255.240.0.0": "12", "255.255.240.0": "20", "255.255.255.240": "28",
		"248.0.0.0": "5", "255.248.0.0": "13", "255.255.248.0": "21", "255.255.255.248": "29",
		"225.0.0.0": "6", "255.252.0.0": "14", "255.255.252.0": "22", "255.255.255.252": "30",
		"254.0.0.0": "7", "255.254.0.0": "15", "255.255.254.0": "23", "255.255.255.254": "31",
		"255.0.0.0": "8", "255.255.0.0": "16", "255.255.255.0": "24", "255.255.255.255": "32",
	}
)

type ipz struct {
	IP      net.IP
	IPNet   *net.IPNet
	IPrange IPRange
}

func NewIPz(args []string) (ipz, error) {
	var ipnet string
	switch len(args) {
	case 1:
		ipnet = args[0]
		i := strings.Split(args[0], "/")
		if len(i) == 2 && len(i[1]) > 2 {
			netmask := rmask_mod[i[1]]
			ipnet = i[0] + "/" + netmask
		} else {
			return ipz{}, fmt.Errorf("输入的参数数量错误，请确认后再试！")
		}
	case 2:
		netmask := args[1]
		if len(args[1]) > 2 {
			netmask = rmask_mod[args[1]]
		}
		ipnet = args[0] + "/" + netmask
	default:
		return ipz{}, fmt.Errorf("输入的参数数量错误，请确认后再试！")
	}
	ip, ipNet, err := net.ParseCIDR(ipnet)
	if err != nil {
		return ipz{}, fmt.Errorf("输入的IP信息有误，请确认后再试！")
	}
	return ipz{
		IP:      ip,
		IPNet:   ipNet,
		IPrange: NewIPRange(ipnet),
	}, nil
}

func bytesToBinary(i []byte) (binary uint32) {
	binaryString := biu.BytesToBinaryString(i)
	biu.ReadBinaryString(binaryString, &binary)
	return
}

func binaryToIP(ipBin uint32) net.IP {
	binaryString := biu.ToBinaryString(ipBin)
	split := strings.Split(binaryString, " ")
	var a, b, c, d uint8
	biu.ReadBinaryString(split[0], &a)
	biu.ReadBinaryString(split[1], &b)
	biu.ReadBinaryString(split[2], &c)
	biu.ReadBinaryString(split[3], &d)
	ipString := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	ip := net.ParseIP(ipString)
	return ip
}

func (ipz *ipz) Run() {
	mask := binaryToIP(bytesToBinary(ipz.IPNet.Mask)).String()
	fmt.Printf("输入的IP地址： %v\n", ipz.IP)
	fmt.Printf("子网掩码信息： %v\n", mask)
	fmt.Printf("所属子网信息： %v\n", ipz.IPNet)
	fmt.Printf("所属子网编号： %v\n", ipz.IPrange.IPID)
	fmt.Printf("所属子网广播： %v\n", ipz.IPrange.IPBro)
	fmt.Printf("主机数量信息： %v\n", ipz.IPrange.hostNumber)
}

func (ipz *ipz) Ranger() {
	out := make(chan string)
	ipz.IPrange.Output(out)
	var i uint32 = 0
	fmt.Println("=======IP地址可用范围=======")
	for ; i < ipz.IPrange.hostNumber; i++ {
		fmt.Println(<-out)
	}
}

type IPRange struct {
	IPID       net.IP
	IPBro      net.IP
	hostNumber uint32
}

func NewIPRange(ipnet string) IPRange {
	_, ipNet, err := net.ParseCIDR(ipnet)
	if err != nil {
		log.Println("输入IP地址段信息有误，请确认后重新尝试！")
	}
	// 获取子网号和广播号
	startIP := bytesToBinary(ipNet.IP)
	mask := bytesToBinary(ipNet.Mask)
	endIP := startIP + ^mask
	// 获取主机数量
	hostNumber := 4294967295 - mask - 1
	return IPRange{
		IPID:       binaryToIP(startIP),
		IPBro:      binaryToIP(endIP),
		hostNumber: hostNumber,
	}
}

func (ipr IPRange) Output(out chan string) {
	go func() {
		ipbin := bytesToBinary(ipr.IPID[12:])
		for {
			ipbin++
			ip := binaryToIP(ipbin)
			out <- ip.String()
		}
	}()
}

func (ipr IPRange) Len() uint32 {
	return ipr.hostNumber
}
