package main

import (
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"github.com/rshulabs/Zinx_/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("call router handle")
	fmt.Println("recv from client : msgid = ", req.GetMsgId(), " data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(1, []byte("ping..."))
	if err != nil {
		fmt.Println("call router handle err:", err)
	}
}

func main() {
	s := znet.NewServer("[myZinx v0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
