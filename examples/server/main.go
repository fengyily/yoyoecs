package main

import (
	"fmt"
	"sync"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protoc"
	"github.com/fengyily/yoyoecs/protocols"
	"google.golang.org/protobuf/proto"
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
			info := protoc.Register{}
			proto.Unmarshal(data, &info)

			server.AddClient(cs.ConnectId, cs)
			cs.SendMessage(protocols.RESPONSE_REGISTER_SUCCESS, 0, []byte(fmt.Sprintf("%v|%v|%v 你注册成功了。", info.ShopCode, info.IP, info.CompanyID)))
		}
		if header.Cmd == protocols.REQUEST_UPLOAD_SKU_DATA {
			if (protocols.Flag(header.Flag) & protocols.HEADER_FLAG_DATA_TYPE_PB) > 0 {

				list := protoc.SkuList{}
				proto.Unmarshal(data, &list)

				fmt.Println("收到了sku数据", header.Length, cs.ConnectId, list.GetSku())
			}
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
