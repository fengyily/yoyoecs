package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protocols"
	"github.com/fengyily/yoyoecs/utils"
)

var sendLock sync.Mutex

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
		if header.Cmd == protocols.RESPONSE_REGISTER_SUCCESS {
			fmt.Println("收到注册成功消息")

			test(cs.GetConn())
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

func test(client *net.Conn) {
	
	go func(client *net.Conn) {
		
		index := 1
		n := 0
		for i := 1; i <= 100; i++ {
			for j := 1; j <= 100; j++ {
				
				index++
				SendMessage(client, protocols.RESPONSE_TRANS_SKU_DATA, []byte("「"+strconv.Itoa(index)+"」这是ＳＫＵ信息，第"+strconv.Itoa(i)+" 轮测试"))
				if n%20 == 0 {
					index++
					SendSplitMessage(client, protocols.REQUEST_UPLOAD_SKU_DATA, []byte("「"+strconv.Itoa(index)+"」这是ＳＫＵ信息，第"+strconv.Itoa(i)+" 批分裂消息"))
				}
				if n%100 == 0 {
					index++
					SendBadMessage(client, protocols.REQUEST_UPLOAD_SKU_DATA, []byte("「"+strconv.Itoa(index)+"」这是一条坏消息，它将影响下一条消息。就看你的程序会不会识别，不要影响下下条消息。。。。"))
				}
				n++
			}
			time.Sleep(2 * time.Second)
		}
	}(client)

}

//正常包发送测试
func SendMessage(conn *net.Conn, cmd protocols.Command, body []byte) (err error) {
	defer func(){
		recover()
	}()
	sendLock.Lock()
	defer sendLock.Unlock()
	data := make([]byte, 2)
	data[0] = byte(cmd)
	data[1] = 2
	(*conn).Write(data)
	if body != nil {
		(*conn).Write(utils.Uint16ToBytes(uint16(len(body))))

		(*conn).Write(body)
	}

	//fmt.Println("test", err, n)
	return
}

// 分裂消息测试
func SendSplitMessage(conn *net.Conn, cmd protocols.Command, body []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sendLock.Lock()
	defer sendLock.Unlock()
	data := make([]byte, 2)
	data[0] = byte(cmd)
	data[1] = 2
	(*conn).Write(data)
	if body != nil {
		(*conn).Write(utils.Uint16ToBytes(uint16(len(body))))

		(*conn).Write(body[0 : len(body)-12])
		time.Sleep(time.Second * 1)
		(*conn).Write(body[len(body)-12:])
	}
	return
}

//　坏包发送测试
func SendBadMessage(conn *net.Conn, cmd protocols.Command, body []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sendLock.Lock()
	defer sendLock.Unlock()
	data := make([]byte, 2)
	data[0] = byte(cmd)
	data[1] = 2
	(*conn).Write(data)
	if body != nil {
		(*conn).Write(utils.Uint16ToBytes(uint16(len(body))))

		(*conn).Write(body[0 : len(body)-5])
	}
	return
}
