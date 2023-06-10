package znet

import (
	"fmt"
	"github.com/rshulabs/Zinx_/utils"
	"github.com/rshulabs/Zinx_/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (h *MsgHandle) DoMsgHandler(req ziface.IRequest) {
	handle, ok := h.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("api msgid = ", req.GetMsgId(), " is not found")
		return
	}
	handle.PreHandle(req)
	handle.Handle(req)
	handle.PostHandle(req)
}

func (h *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := h.Apis[msgId]; ok {
		panic("repeated api msgid = " + strconv.Itoa(int(msgId)))
	}
	h.Apis[msgId] = router
	fmt.Println("add api msgid = ", msgId)
}

func (h *MsgHandle) StartOneWorker(workID int, queue chan ziface.IRequest) {
	fmt.Println("worker id = ", workID, " is started.")
	for {
		select {
		case req := <-queue:
			h.DoMsgHandler(req)
		}
	}
}

func (h *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		h.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go h.StartOneWorker(i, h.TaskQueue[i])
	}
}

func (h *MsgHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	workID := req.GetConnection().GetConnID() % h.WorkerPoolSize
	fmt.Println("add connid = ", req.GetConnection().GetConnID(), " req msgid = ", req.GetMsgId(), " to workid = ", workID)
	h.TaskQueue[workID] <- req
}
