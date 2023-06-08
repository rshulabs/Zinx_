package ziface

/**
len id data
*/

type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte
	SetData([]byte)
	SetDataLen(uint32)
}
