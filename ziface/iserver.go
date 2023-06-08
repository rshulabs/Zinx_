package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// AddRouter 添加路由
	AddRouter(router IRouter)
}
