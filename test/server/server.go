package main

import (
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"github.com/rshulabs/Zinx_/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("call router prehandle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back prehandle err:", err)
	}
}

func (r *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("call router handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call router handle err:", err)
	}
}

func (r *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("call router posthandle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back posthandle err:", err)
	}
}

func main() {
	s := znet.NewServer("[myZinx v0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
