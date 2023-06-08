package main

import (
	"fmt"
	"github.com/rshulabs/Zinx_/znet"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client conn err:", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("zinx_ v0.4")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client send data err:", err)
			return
		}
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client recv data err:", err)
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("recv msg: id = ", msg.Id, " len = ", msg.DataLen, " data = ", string(msg.Data))
		}
		time.Sleep(time.Second)
	}
}
