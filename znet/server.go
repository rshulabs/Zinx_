package znet

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func CallbackToClient(conn *net.TCPConn, buf []byte, cnt int) error {
	fmt.Println("[connection handle] CallbackToClient")
	if _, err := conn.Write(buf[:cnt]); err != nil {
		fmt.Println("send data err:", err)
		return errors.New("callback err")
	}
	return nil
}
func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP : %s,Port : %d\n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcpaddr err:", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listener err:", err)
			return
		}
		fmt.Println("start server zinx ", s.Name, " success, now listening...")
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}
			// TODO

			// TODO
			dealConn := NewConnection(conn, cid, CallbackToClient)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server,name:", s.Name)
	// TODO
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
}
