package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// AddRouter 添加路由
	AddRouter(msgId uint32, router IRouter)
}
