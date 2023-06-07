package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client conn err:", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello i am client"))
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
		fmt.Printf("server call back : %s,cnt : %d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}

func TestServer(t *testing.T) {
	s := NewServer("[myZinx v0.1]")
	go Client()
	s.Serve()
}
