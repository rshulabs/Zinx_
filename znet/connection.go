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
	handleAPI ziface.HandFunc
	// 告知该连接已经退出的channel
	ExistBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, id uint32, callback ziface.HandFunc) *Connection {
	return &Connection{
		Conn:          conn,
		ConnID:        id,
		handleAPI:     callback,
		isClosed:      false,
		ExistBuffChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutinue is running")
	defer fmt.Println(c.getRemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			c.ExistBuffChan <- true
			continue
		}
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID", c.ConnID, " handle err", err)
			c.ExistBuffChan <- true
			return
		}
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

func (c *Connection) getTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) getConnID() uint32 {
	return c.ConnID
}

func (c *Connection) getRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
