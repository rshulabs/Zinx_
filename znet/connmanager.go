package znet

import (
	"errors"
	"fmt"
	"github.com/rshulabs/Zinx_/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (m *ConnManager) Add(conn ziface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()
	m.connections[conn.GetConnID()] = conn
	fmt.Println("conn add to manager success , len =  ", m.Len())
}

func (m *ConnManager) Remove(conn ziface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()
	delete(m.connections, conn.GetConnID())
	fmt.Println("conn remove id = ", conn.GetConnID())
}

func (m *ConnManager) Get(id uint32) (ziface.IConnection, error) {
	m.connLock.RLock()
	defer m.connLock.RUnlock()
	if conn, ok := m.connections[id]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn not found")
	}
}

func (m *ConnManager) Len() int {
	return len(m.connections)
}

func (m *ConnManager) ClearConn() {
	m.connLock.Lock()
	defer m.connLock.Unlock()
	for id, conn := range m.connections {
		conn.Stop()
		delete(m.connections, id)
	}
	fmt.Println("clear all conn success len = ", m.Len())
}
