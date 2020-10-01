package trace_cli

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
)

func Tracesite(ctx *context) error {

	// Address Family: AF,地址族，ip地址这边（特指V4）
	domain := syscall.AF_INET
	// 协议族 protocol family。比如udp和tcp。
	proto_send := syscall.IPPROTO_UDP
	proto_recv := syscall.IPPROTO_ICMP
	// 核心是调用syscall的socket。这个函数需要哪些选项就是context要囊括进去
	// udp要设置datagram类型包
	send_socket, err := syscall.Socket(domain, syscall.SOCK_DGRAM, proto_send)
	if err != nil {
		return err
	}
	// 使用raw-socket是比传输层更底层得靠用户自己提供和设置的报文，可以方便获取报文内容
	recv_socket, err := syscall.Socket(domain, syscall.SOCK_RAW, proto_recv)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer syscall.Close(send_socket)
	defer syscall.Close(recv_socket)

	ttl := ctx.begin
	timeout := syscall.NsecToTimeval(1000 * 1000 * int64(ctx.timeout))
	port := ctx.port
	// 最终传进去的地址都要转化为字节（数组）
	srcAddrBytes, err := generateAddrInBytes()
	dstAddrBytes, addrs, err := generateDstAddrInBytes(ctx.addr)
	retries := 0
	dstStr := fmt.Sprintf("%v.%v.%v.%v", dstAddrBytes[0],
		dstAddrBytes[1], dstAddrBytes[2], dstAddrBytes[3])
	fmt.Printf("Traceroute to %v (%v), %d hops max, 64 byte packets\n",
		ctx.addr, dstStr, ctx.maxhop)
	first := true
	for {
		if ttl > ctx.maxhop {
			fmt.Println("ttl > ctx.maxhop")
			break
		}

		// 这些选项设置的解释可以参考unix网络编程那本书
		// opt int是通过int来设置对应的选项。opt是选项ID，value是选项值
		syscall.SetsockoptInt(send_socket, 0x00,
			syscall.IP_TTL, ttl)
		// 设置接收超时
		syscall.SetsockoptTimeval(recv_socket, syscall.SOL_SOCKET,
			syscall.SO_RCVTIMEO, &timeout)

		// sa source addr
		source_addr := &syscall.SockaddrInet4{
			Port: port,
			Addr: srcAddrBytes,
		}
		// 接收监听端口
		syscall.Bind(recv_socket, source_addr)
		dst_addr := &syscall.SockaddrInet4{
			Port: port,
			Addr: dstAddrBytes,
		}
		// 发送的目的socket
		syscall.Sendto(send_socket, []byte{0x00}, 0, dst_addr)

		startTime := time.Now()
		// packet-size 包的大小（单位字节）
		p := make([]byte, ctx.packetsize)
		// 接收，就阻塞在这边。n为读取的大小，from不一定是dst，from这边表示
		// 返回数据包的地址
		_, from, err :=syscall.Recvfrom(recv_socket, p, 0)

		elapseTime := time.Since(startTime)
		if err == nil {
			retries = 0
			ip := from.(*syscall.SockaddrInet4).Addr
			fromAddrStr := fmt.Sprintf("%v.%v.%v.%v",
				ip[0], ip[1], ip[2], ip[3])
			// 这边符合一个即可退出
			if addrs[fromAddrStr]{
				fmt.Println("addrs[fromAddrStr]这边退出")
				break
			}
			if first {
				fmt.Printf("%v. %v // %v\n", ttl,
					fromAddrStr, elapseTime)
			}else {
				fmt.Printf("%v // %v\n",
					fromAddrStr, elapseTime)
			}
			first = true
			ttl ++
		}else {
			// 检查重试了几次了已经
			if retries < ctx.retry {
				first = false
				if retries == 0 {
					fmt.Printf("%v. ", ttl)
				}
				retries ++
				fmt.Printf("* ")

			}else {
				first = true
				retries = 0
				// 过不去某个节点了，去下一个节点试试看
				fmt.Printf("当前节点一直不通，将跳过当前节点\n")
				ttl += 1
			}
		}
	}
	return nil
}



// 涉及到net包的学习
// 选择本地可用的单播地址作为发送socket
func generateAddrInBytes() ([4]byte, error)  {
	socketAddr := [4]byte{0, 0, 0, 0}
	// 获取本地IP
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return socketAddr, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// To4其实就是转化为字节数组
			if len(ipnet.IP.To4()) == net.IPv4len {
				// 数组转变为slice的最直接方式[:]
				fmt.Println("可用的原地址: ", ipnet.IP.String())
				copy(socketAddr[:], ipnet.IP.To4())
				return socketAddr, nil
			}
		}
	}
	return socketAddr, errors.New("本机未联网")
}

// 涉及到net包的学习
func generateDstAddrInBytes(dst string) ([4]byte, map[string]bool, error)  {
	dstAddr := [4]byte{0, 0, 0, 0}
	// 将域名转化为字符串再转化为字节数组
	// 下面这行代码先转为字符串数组，因为一个域名可能在本机存储了多个IP地址和
	addrs, err := net.LookupHost(dst)
	// 这边可能返回了多个可用addrs，
	if err != nil {
		return dstAddr, nil,  err
	}
	// ip地址字符串转化为封装了IP地址的结构体，该结构体有转化为字节数组的函数
	// 优先使用第一个addr作为目标地址
	ipaddr, err := net.ResolveIPAddr("ip", addrs[0])
	if err != nil {
		return dstAddr, nil, err
	}

	addrMap := map[string]bool{}
	for _, elem := range addrs{
		addrMap[elem] = true
	}
	copy(dstAddr[:], ipaddr.IP.To4())
	return dstAddr, addrMap, nil
}