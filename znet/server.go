package znet

import (
	"fmt"
	"github.com/rshulabs/Zinx_/utils"
	"github.com/rshulabs/Zinx_/ziface"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	// 由用户绑定回调router
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP : %s,Port : %d\n", s.IP, s.Port)
	fmt.Printf("[ZINX_] version: %s, maxconn: %d, macpacketsize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
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
			dealConn := NewConnection(conn, cid, s.Router)
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
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("add router success")
}

func NewServer(name string) *Server {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
