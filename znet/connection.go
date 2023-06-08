package znet

import (
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"net"
)

type Connection struct {
	// 当前连接socket
	Conn *net.TCPConn
	// 当前连接id
	ConnID uint32
	// 当前连接关闭状态
	isClosed bool
	// 处理api
	//handleAPI ziface.HandFunc

	// 处理router
	Router ziface.IRouter
	// 告知该连接已经退出的channel
	ExistBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, id uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:   conn,
		ConnID: id,
		//handleAPI:     callback,
		Router:        router,
		isClosed:      false,
		ExistBuffChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutinue is running")
	defer fmt.Println(c.GetRemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			c.ExistBuffChan <- true
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		go func(req ziface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <-c.ExistBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	c.ExistBuffChan <- true
	close(c.ExistBuffChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
