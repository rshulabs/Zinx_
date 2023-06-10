package utils

import (
	"encoding/json"
	"github.com/rshulabs/Zinx_/ziface"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer        ziface.IServer `json:"tcp_server"`
	Host             string         `json:"host"`
	TcpPort          int            `json:"tcp_port"`
	Name             string         `json:"name"`
	Version          string         `json:"version"`
	MaxPacketSize    uint32         `json:"max_packet_size"` // 允许最大包大小
	MaxConn          int            `json:"max_conn"`        // 允许最大连接数
	WorkerPoolSize   uint32         `json:"worker_pool_size"`
	MaxWorkerTaskLen uint32         `json:"max_worker_task_len"`

	ConfFilePath string `json:"conf_file_path"`
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("../../conf/zinx_.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 饿汗
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TcpPort:          7777,
		Host:             "127.0.0.1",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "../conf/zinx_.json",
	}
	// 从配置读取
	GlobalObject.Reload()
}
