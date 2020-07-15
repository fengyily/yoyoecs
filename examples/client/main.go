package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protocols"
)

func main() {
	client := yoyoecs.ClientSocket{}
	client.OnConnError = func(err error) {
		fmt.Println("on connect error ", err)
	}
	client.OnRecvMessage = func(header protocols.Header, data []byte, cs *yoyoecs.ClientSocket) {
		fmt.Println("-----------------------------------------------------------", time.Now())
		fmt.Println("成功收到消息：", header.Cmd, "长度：", header.Length, "Flag", fmt.Sprintf("%08b", header.Flag), string(data))
		fmt.Println("-----------------------------------------------------------", time.Now())

		if header.Cmd == protocols.RESPONSE_UPLOAD_SKU_DATA {
			fmt.Println("通知继续发送。")
		}
		if (header.Cmd == protocols.RESPONSE_REGISTER_SUCCESS) {
			fmt.Println("收到注册成功消息")
		}
	}

	client.OnConnect = func(ip string, cs *yoyoecs.ClientSocket) {
		type EdgeRegister struct {
			IP        string `json:"ip"`
			SN        string `json:"sn"`
			CompanyID int64  `json:"company_id"`
			ShopCode  string `json:"shop_code"`
		}
		info := EdgeRegister{}
		info.CompanyID = 123456789
		info.ShopCode = "shopcode123456789"
		d, _ := json.Marshal(info)
		client.SendMessage(protocols.REQUEST_REGISTER, 0, d)
		fmt.Println("发起了注册申请")
	}
	client.Conn("127.0.0.1:9091")

	sw := sync.WaitGroup{}
	sw.Add(1)
	sw.Wait()
}
