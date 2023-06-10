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
	msgHandler ziface.IMsgHandle
	// 连接
	ConnMgr ziface.IConnManager

	// Hook
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP : %s,Port : %d\n", s.IP, s.Port)
	fmt.Printf("[ZINX_] version: %s, maxconn: %d, macpacketsize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	go func() {
		s.msgHandler.StartWorkerPool()
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
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			// TODO
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server,name:", s.Name)
	// TODO
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	select {}
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("add router success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hook func(connection ziface.IConnection)) {
	s.OnConnStart = hook
}

func (s *Server) SetOnConnStop(hook func(connection ziface.IConnection)) {
	s.OnConnStop = hook
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---->callonconnstart...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---->callonconnstop...")
		s.OnConnStop(conn)
	}
}