package ziface

import "net"

type IConnection interface {
	// Start 开启连接
	Start()
	// Stop 关闭连接
	Stop()
	// GetTcpConnection 获取当前连接socket
	GetTcpConnection() *net.TCPConn
	// GetConnID 获取当前连接id
	GetConnID() uint32
	// GetRemoteAddr 获取客户端地址
	GetRemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error
	SendBuffMsg(msgId uint32, data []byte) error
}

// HandFunc 定义一个统一处理业务接口
//type HandFunc func(*net.TCPConn, []byte, int) error
