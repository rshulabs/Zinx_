package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 立刻以非阻塞处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体处理逻辑
}
