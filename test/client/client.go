package main

import (
	"fmt"
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
		_, err := conn.Write([]byte("hello Zinx"))
		if err != nil {
			fmt.Println("client send data err:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client recv data err:", err)
			return
		}
		fmt.Printf("server call back : %s, cnt : %d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}
