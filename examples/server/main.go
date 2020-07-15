package main

import (
	"fmt"
	"sync"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protocols"
)

func main() {
	server := &yoyoecs.ServerSocket{}
	server.OnConnect = func(ip string, cs *yoyoecs.ClientSocket) {
		fmt.Println(ip, "客户端连接成功了")
	}
	server.OnError = func(cs *yoyoecs.ClientSocket) {
		fmt.Println("连接断开，把连接移除之，", cs.ConnectId)
		server.RemoveClient(cs.ConnectId)
	}
	server.OnRecvMessage = func(header protocols.Header, data []byte, cs *yoyoecs.ClientSocket) {
		fmt.Println("收到消息", header.Cmd, "Length", header.Length)
		if header.Cmd == protocols.REQUEST_HEARTBEAT {
			server.ClientOnline(cs.ConnectId, cs)
		}
		if header.Cmd == protocols.REQUEST_REGISTER {
			server.AddClient(cs.ConnectId, cs)
			cs.SendMessage(protocols.RESPONSE_REGISTER_SUCCESS, 0, []byte("你注册成功了。"))
		}
		if header.Cmd == protocols.REQUEST_UPLOAD_SKU_DATA {
			fmt.Println("收到了sku数据", header.Length, cs.ConnectId)
			cs.SendMessage(protocols.RESPONSE_PASSIVE_UPLOAD_SKU_DATA, 0, []byte("收到，请继续发送，如果你还有的话。"))
		}
	}

	server.Run(":9091")
	fmt.Println("start success, wait for you connect now.")
	//<-time.After(10*time.Second)

	//server.Close()

	sw := sync.WaitGroup{}
	sw.Add(1)
	sw.Wait()
}
