package trace_cli

import (
	"flag"
	"fmt"
)

type context struct {
	begin      	int
	maxhop		int
	timeout    	int
	port       	int
	addr       	string
	packetsize 	int
	retry      	int
}

type ActionFunc func(*context) error

type trace_cli struct {

	Action ActionFunc
	Name string
	Usage string
	UsageHelp string
}
const DEFAULT_PORT = 34456
const DEFAULT_PACKET_SIZE = 64
const DEFAULT_ADDR = "www.baidu.com"
const DEFAULT_RETRIES = 3
const DEFAULT_MAX_HOP = 64
const DEFAULT_TIMEOUT = 3000
const DEFAULT_BEGIN = 1
var timeout int
// 从哪个跳数开始
var begin int
// 最大跳数
var maxhop int
var port int
func (c *trace_cli) run(args []string) error {
	flag.IntVar(&timeout, "timeout", DEFAULT_TIMEOUT, "设置超时")
	flag.IntVar(&begin, "begin", DEFAULT_BEGIN, "设置起点的跳数")
	flag.IntVar(&port, "p", DEFAULT_PORT, "自定义端口")
	flag.IntVar(&maxhop, "maxhop", DEFAULT_MAX_HOP, "自定义最大跳数")
	flag.Parse()
	ctx := &context{
		begin:      begin,
		timeout:    timeout,
		port:       port,
		addr:       DEFAULT_ADDR,
		packetsize: DEFAULT_PACKET_SIZE,
		retry:      DEFAULT_RETRIES,
		maxhop: 	maxhop,
	}
	fmt.Println("超时时间是: ", timeout)
	// 解析args参数，创建上下文(其实就是网络请求需要的一些option)
	// 有了flag还要os.Args是因为用户的输入是不确定的。如果没找到网址(网址可以是直接IP地址，也可以是域名)
	// 这边一样要返回错误。
	ctx.addr = args[len(args) - 1]


	// 执行action
	err := c.Action(ctx)
	return err
}

func Run(args []string) error {
	// 创建客户端
	app := &trace_cli{
		Action: func(ctx *context) error {
			if len(args) >= 2 {
				return Tracesite(ctx)
			}else {
				return fmt.Errorf("参数长度小于等于1，没有指定地址。输入--help查看使用说明。")
			}
		},
		Usage: "这是一个用Go实现的路由追踪客户端",
		UsageHelp: "trace --begin=3 --timeout=2000 www.anywebsite.com\t网站名字要放在最后",
		Name: "路由追踪应用",
	}

	return app.run(args)
}