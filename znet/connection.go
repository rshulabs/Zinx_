package znet

import (
	"errors"
	"fmt"
	"github.com/rshulabs/Zinx_/utils"
	"github.com/rshulabs/Zinx_/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TcpServer ziface.IServer
	// 当前连接socket
	Conn *net.TCPConn
	// 当前连接id
	ConnID uint32
	// 当前连接关闭状态
	isClosed bool
	// 处理api
	//handleAPI ziface.HandFunc

	// 处理消息
	MsgHandler ziface.IMsgHandle
	// 告知该连接已经退出的channel
	ExistBuffChan chan bool

	// 读写goroutinue通信
	msgChan chan []byte
	// 读写goroutinue通信 有缓冲
	msgBuffChan chan []byte

	// 属性
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, id uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:     server,
		Conn:          conn,
		ConnID:        id,
		MsgHandler:    msgHandler,
		isClosed:      false,
		ExistBuffChan: make(chan bool, 1),
		msgChan:       make(chan []byte),
		msgBuffChan:   make(chan []byte, utils.GlobalObject.MaxWorkerTaskLen),
		property:      make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
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
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.ExistBuffChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExistBuffChan)
	close(c.msgBuffChan)
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
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send buff data err: ", err)
					return
				}
			} else {
				fmt.Println("msgbuff is closed")
				break
			}
		case <-c.ExistBuffChan:
			return

		}
	}
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("conn closed when send buf msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack err msg id = ", msgId)
		return errors.New("pack err msg")
	}
	c.msgBuffChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = val
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if val, ok := c.property[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
