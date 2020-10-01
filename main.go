package main

import (
	c "github.com/Ichliebedich0505/T_traceroute/trace_cli"
	"log"
	"os"
)

// 程序本身就是command
var timeout string

const defaultSite = "https://www.baidu.com"
func main() {
	// 用flag来解析命令在这边不太好用，已经废弃
	//// 接收命令。要先声明才调用parse
	//flag.StringVar(&timeout, "timeout", defaultSite,
	//	"跟踪一个网址：trace --timeout=2000 网址")
	//flag.Parse()
	//
	//// 解析出命令的option和要追踪的网站
	//fmt.Println("timeout: ", timeout)
	//fmt.Println("os.args: ", os.Args)
	if err := c.Run(os.Args); err != nil {
		log.Printf("%v", err)
	}
}
