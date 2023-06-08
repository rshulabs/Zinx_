package ziface

type IRouter interface {
	PreHandle(request IRequest) // 处理连接之前
	Handle(request IRequest)
	PostHandle(request IRequest) // 处理连接之后
}
