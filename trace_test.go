package main

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

type a struct {
	i int
}

type b struct {
	a
	l int
}

func TestArgs(t *testing.T) {
	//fmt.Println(os.Args)
	addrs, err := net.LookupHost("www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(addrs)

	ipaddr, err := net.ResolveIPAddr("ip", addrs[0])
	fmt.Println(err)
	fmt.Println(ipaddr)
}


func TestInterface(t *testing.T) {
	bb := b{
		a: a{6},
		l: 1,
	}
	fmt.Println(bb)
	app(3)
}

func app(a ...interface{}) {

	for aa := range a {
		fmt.Println(reflect.TypeOf(aa))
	}

}
