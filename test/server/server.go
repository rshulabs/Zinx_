package main

import (
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"github.com/rshulabs/Zinx_/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

type HelloRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("call pinghandle")
	fmt.Println("recv from client : msgid = ", req.GetMsgId(), " data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(0, []byte("ping..."))
	if err != nil {
		fmt.Println("call pinghandle err:", err)
	}
}

func (r *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("call hellohandle")
	fmt.Println("recv from client : msgid = ", req.GetMsgId(), " data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(1, []byte("hello world!"))
	if err != nil {
		fmt.Println("call hellohandle err:", err)
	}
}
func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
