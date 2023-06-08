package znet

import (
	"errors"
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"io"
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
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTcpConnection(), headData)
		if err != nil {
			fmt.Println("recv msg len err:", err)
			c.ExistBuffChan <- true
			continue
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			c.ExistBuffChan <- true
			continue
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("recv msg data err:", err)
				c.ExistBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack err msg id = ", msgId)
		return errors.New("pack err msg")
	}
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("write msg id ", msgId, " err")
		c.ExistBuffChan <- true
		return errors.New("conn write err")
	}
	return nil
}
