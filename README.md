# mytools
之前也用过python和shell做过一些程序比如IP地址计算器、多ping小程序，而这次使用go语言实现这些功能。
直接克隆到本地，使用go install直接可以获得执行程序。
## mping
多ping小程序网上一搜一大把，但是都不够快捷。我主要将**github.com/sparrc/go-ping**做了二次封装集成了多IP地址，多次PING测试。
并且可以把文件或者IP地址段作为输入进行测试操作。在将地址段作为输入的时候也不必输入子网号/子网掩码，只要将子网中IP/子网掩码输入即可。
而文件输入只要将IP地址依照换行进行存放即可。
```shell script
> mytools help mping
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  mytools mping [flags]

Flags:
  -c, --count int         每次每个对象测试包数 (default 5)
  -h, --help              help for mping
  -i, --interval int      每个测试对象发包的时间间隔，单位毫秒 (default 500)
  -f, --ipfile string     从文件中获取测试IP地址
  -n, --ipnet string      根据网段获取测试IP地址，例如192.168.0.0/24
  -r, --repeat int        所有IP测试的遍数，默认1次 (default 1)
  -s, --spacing int       每次测试时间间隔 (default 2)
  -w, --workercount int   同时测试的协程数量 (default 4)
> mytools mping -n 192.168.0.1/30
==================第1次测试=================
2020/07/05 15:15:09 192.168.0.1 ping statistics: 5 packets transmitted, 5 received, 0% packet loss, max 2.0011ms
2020/07/05 15:15:15 192.168.0.4 ping statistics: 5 packets transmitted, 0 received, 100% packet loss, max 0s
```

## ipz
这个模块主要是IPv4地址计算器，可以输入192.168.0.1/24，也可以输入192.168.0.1/255.255.255.252，
还可以输入192.168.0.1 255.255.255.252。使用-r还可以得到可用主机信息。
```shell script
> go run main.go ipz 192.168.0.1/255.255.255.252 -r
输入的IP地址： 192.168.0.1
子网掩码信息： 255.255.255.252
所属子网信息： 192.168.0.0/30
所属子网编号： 192.168.0.0
所属子网广播： 192.168.0.3
主机数量信息： 2
=======IP地址可用范围=======
192.168.0.1
192.168.0.2
```

## recactibw
这个模块主要是修改cacti监控交换机端口进出流量的rrd文件。主要是因为对账的同事要克扣供应商、抬高客户计费（这里强烈谴责），另外就是修改较大的突发
流量、补漏数据、查询数据等功能。其中原理主要是将rrd文件转换为xml格式进行修改，所以xml中的row数据时间间隔必须符合rrd文件中step的数值。如有不对
将会导致数据错乱，使用rrdtool dump xxx.rrd > xxx.xml导出到文件中。查看<step> 300 </step> <!-- Seconds -->中为300秒，database中也
满足row间隔300秒。
```xml
<!-- 2019-08-05 14:35:00 CST / 1564986900 --> <row><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v></row>
<!-- 2019-08-05 14:40:00 CST / 1564987200 --> <row><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v></row>
<!-- 2019-08-05 14:45:00 CST / 1564987500 --> <row><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v></row>
<!-- 2019-08-05 14:50:00 CST / 1564987800 --> <row><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v><v> NaN </v></row>
```
命令提示如下
```shell
Flags:
  -d, --details           是否打印详情
      --end string        修改起始时间点2020-01-21 00:00:00
  -h, --help              help for recactibw
  -i, --interval int      参照若干时间点前的数据修补当前数据单位为天
  -l, --log string        操作文件日志所存地方 (default "/var/log/mytools/recactibw.log")
  -p, --percent float     将一定范围内的数据按照百分百升降.0.3(表示减小到30%) (default 1)
      --print             显示数值模式
  -s, --source string     需要更新的目标文件或文件夹
      --start string      修改起始时间点2020-01-21 00:00:00
  -w, --workercount int   同时运行修改文件的数量 (default 5)
go run main.go recactibw  -i 7 -p 1.1 --start "2021-01-12 05:00:00" --end "2021-01-17 23:59:00" -w 7 -s /var/www/html/rra/
```
修改start到end时间点内的数据，-i 7 每个时间点参照7天前相同时间段的数据，-p 1.1 乘以1.1，-w 7同时运行7个协程。-s数据源可以是文件或文件夹