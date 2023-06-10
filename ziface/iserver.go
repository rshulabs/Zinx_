package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// AddRouter 添加路由
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 连接
	GetConnMgr() IConnManager

	// Hook
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}
