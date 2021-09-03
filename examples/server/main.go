/*
 * @Author: F1
 * @Date: 2020-07-15 09:36:41
 * @LastEditors: F1
 * @LastEditTime: 2021-09-02 22:26:37
 * @Description: 服务端的测示例
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protoc"
	"github.com/fengyily/yoyoecs/protocols"
	"google.golang.org/protobuf/proto"
)

func main() {
	server := &yoyoecs.ServerSocket{}
	server.OnConnect = func(ip string, cs *yoyoecs.ClientSocket) {
		fmt.Println(ip, "客户端连接成功了", cs.ConnectId)

	}
	server.OnError = func(cs *yoyoecs.ClientSocket) {
		fmt.Println("连接断开，把连接移除之，", cs.ConnectId)
		server.RemoveClient(cs.ConnectId)
	}
	server.OnSendToMessage = func(serverInfo string, header protocols.Header, data []byte, cs *yoyoecs.ClientSocket) {
		fmt.Println(serverInfo, "收到 OnSendToMessage 消息")
	}
	server.OnRecvMessage = func(header protocols.Header, data []byte, cs *yoyoecs.ClientSocket) {
		fmt.Println("收到消息", header.Cmd, "Length", header.Length)
		if header.Cmd == protocols.REQUEST_HEARTBEAT {
			server.ClientOnline(cs.ConnectId, cs)
		}
		if header.Cmd == protocols.REQUEST_REGISTER {
			info := protoc.Register{}
			proto.Unmarshal(data, &info)

			server.AddClient(info.SN, cs)
			cs.SendMessage(protocols.RESPONSE_REGISTER_SUCCESS, 0, []byte(fmt.Sprintf("SN:%s %v|%v|%v 你注册成功了。", info.SN, info.ShopCode, info.IP, info.CompanyID)))
			//cs.SendMessage(protocols.RESPONSE_UPLOAD_SKU_DATA, 0, []byte(fmt.Sprintf("收到:%v，请继续发送，如果你还有的话。", 1)))
		}
	}

	server.Run(":9091")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("start success, wait for you connect now.")
EXIT:
	for {
		fmt.Println("wait signal")
		sig := <-sc
		fmt.Printf("获取到信号[%s]", sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	server.Close()
}
