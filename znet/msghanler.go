package znet

import (
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
