/*
 * @Author: F1
 * @Date: 2020-07-15 09:36:41
 * @LastEditors: F1
 * @LastEditTime: 2020-08-07 10:16:19
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
			//cs.SendMessage(protocols.RESPONSE_UPLOAD_SKU_DATA, 0, []byte(fmt.Sprintf("收到:%v，请继续发送，如果你还有的话。", 1)))
			cs.SendMessage(protocols.REQUEST_PASSIVE_UPLOAD_SKU_DATA, 0, []byte("快上传"))
		}
		if header.Cmd == protocols.REQUEST_UPLOAD_SKU_DATA {
			// 以下是protocbuf格式传输
			list := protoc.SkuList{}
			proto.Unmarshal(data, &list)

			//cs.SendMessage(protocols.RESPONSE_REGISTER_SUCCESS, 0, []byte("%v|%v|%v 你注册成功了。"))
			cs.SendMessage(protocols.RESPONSE_UPLOAD_SKU_DATA, 0, []byte("请继续"))
			fmt.Println("收到了sku数据", header.Length, cs.ConnectId, list.GetSku())

		}
		if header.Cmd == protocols.REQUEST_TRANS_YOYOINFO_DATA {
			yoyoList := protoc.YoyoInfoList{}
			proto.Unmarshal(data, &yoyoList)

			fmt.Println("收到了yoyo数据", header.Length, cs.ConnectId, yoyoList.GetYoyoInfo())
			cs.SendMessage(protocols.RESPONSE_TRANS_YOYOINFO_DATA, 0, []byte(fmt.Sprintf("收到:%v，请继续发送，如果你还有的话。", len(yoyoList.GetYoyoInfo()))))
		}
		if header.Cmd == protocols.RESPONSE_PASSIVE_UPLOAD_SKU_DATA {
			fmt.Println("收到了", string(data))
			//cs.SendMessage(protocols.RESPONSE_PASSIVE_UPLOAD_SKU_DATA, 0, []byte("continue"))
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
