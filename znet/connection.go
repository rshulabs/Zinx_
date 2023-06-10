package znet

import (
	"errors"
	"fmt"
	"github.com/rshulabs/Zinx_/utils"
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
	MsgHandler ziface.IMsgHandle
	// 告知该连接已经退出的channel
	ExistBuffChan chan bool

	// 读写goroutinue通信
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, id uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:          conn,
		ConnID:        id,
		MsgHandler:    msgHandler,
		isClosed:      false,
		ExistBuffChan: make(chan bool, 1),
		msgChan:       make(chan []byte),
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
		//fmt.Println("head : ", headData)
		if err != nil {
			fmt.Println("recv msg len err:", err)
			c.ExistBuffChan <- true
			continue
		}
		msg, err := dp.Unpack(headData)
		//fmt.Println("msg ", msg)
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
		if utils.GlobalObject.MaxPacketSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
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
	c.msgChan <- msg
	return nil
}

func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running")
	defer fmt.Println(c.GetRemoteAddr().String(), "[conn writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err: ", err, "conn writer exit")
				return
			}
		case <-c.ExistBuffChan:
			return

		}
	}
}
