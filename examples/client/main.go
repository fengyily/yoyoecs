/*
 * @Author: F1
 * @Date: 2020-07-21 11:47:32
 * @LastEditors: F1
 * @LastEditTime: 2021-09-07 21:30:57
 * @Description: 客户端测试
 */
package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/fengyily/yoyoecs"
	"github.com/fengyily/yoyoecs/protoc"
	"github.com/fengyily/yoyoecs/protocols"
	"github.com/fengyily/yoyoecs/utils"
	"google.golang.org/protobuf/proto"
)

var sendLock sync.Mutex

func main() {
	//TestCompress()
	client := yoyoecs.ClientSocket{}
	client.OnConnError = func(err error) {
		fmt.Println("on connect error ", err)
	}
	client.OnRecvMessage = func(header protocols.Header, data []byte, cs *yoyoecs.ClientSocket) {
		fmt.Println("-----------------------------------------------------------", time.Now())
		fmt.Println("成功收到消息：", header.Cmd, "长度：", header.Length, "Flag", fmt.Sprintf("%08b", header.Flag), string(data), "SN:", cs.ConnectId)
		fmt.Println("-----------------------------------------------------------", time.Now())

		if header.Cmd == protocols.RESPONSE_REGISTER_SUCCESS {
			fmt.Println("收到注册成功消息")
		}
		if header.Cmd == protocols.REQUEST_EXEC_CMD {
			shellCmd := &protoc.ShellCmd{}
			err := proto.Unmarshal(data, shellCmd)
			if err != nil {
				fmt.Println("REQUEST_EXEC_CMD:proto.Unmarshal(data, shellCmd) err:", err.Error())
				return
			}

			fmt.Println(shellCmd.Command)

			output := "already exec command:" + shellCmd.Command
			fmt.Println(shellCmd.Command, "done")
			reply := &protoc.ShellExecReply{}
			reply.Result = output
			body, err := proto.Marshal(reply)
			if err == nil {
				err = cs.SendMessage(protocols.RESPONAE_EXEC_CMD_REPLY, protocols.HEADER_FLAG_DATA_TYPE_PB&protocols.HEADER_FLAG_IS_COMPRESS, body)

				fmt.Printf("Reply success, err:%v %#v %v\r\n", err, reply, body)
			} else {
				fmt.Println("Reply failed, proto.Marshal(reply) err:", err.Error())
			}
		}
	}

	client.OnConnect = func(ip string, cs *yoyoecs.ClientSocket) {
		// 采用ＰＢ上传ＳＫＵ信息的示例
		info := protoc.Register{}
		info.SN = "123"
		info.ShopCode = "456"
		d, _ := proto.Marshal(&info)
		client.SendMessage(protocols.REQUEST_REGISTER, protocols.HEADER_FLAG_IS_COMPRESS|protocols.HEADER_FLAG_DATA_TYPE_PB, d)
		fmt.Println("发起了注册申请")
	}
	client.Conn("127.0.0.1:9091")

	sw := sync.WaitGroup{}
	sw.Add(1)
	sw.Wait()
}

/**
 * @Title:
 * @Description:
 *
 * 					测试各种情况包的发送：
 * 					- 大包（超出缓存区的包），用以测试分包接收是否正常
 *					- 分包（一个完整的包，分两次发送），用以测试并包是否正常
 *					- 异常包（不完整的包），用以测试是否能重新定位包
 *
 * @Author: F1
 * @Date: 2020-07-28 09:30:42
 * @Param:
 * @Return:
 */
func test(client *yoyoecs.ClientSocket) {
}

//正常包发送测试
func SendMessage(conn *net.Conn, cmd protocols.Command, body []byte) (err error) {
	defer func() {
		recover()
	}()
	sendLock.Lock()
	defer sendLock.Unlock()
	data := make([]byte, 2)
	data[0] = byte(cmd)
	data[1] = 0
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
	data[1] = 0
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
	data[1] = 0
	(*conn).Write(data)
	if body != nil {
		(*conn).Write(utils.Uint16ToBytes(uint16(len(body))))

		(*conn).Write(body[0 : len(body)-5])
	}
	return
}

/**
 * @Title:TestCompress　测试包的压缩
 * @Description:
 *
 *					测试包的压缩
 *
 * @Author: F1
 * @Date: 2020-07-28 09:33:46
 * @Param:
 * @Return:
 */
func TestCompress(t *testing.T) {
	befor := []byte("这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的")
	after := utils.Compress(befor)
	println("befor", len(befor))

	println("after", len(after), string(after))

	out := utils.UnCompress(after)

	println("out", string(out))
}

func TestPB(t *testing.T) {
	reg := &protoc.Register{}

	reg.SN = "SN00112343"
	reg.CompanyID = 12345674879
	reg.ShopCode = "这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的"
	reg.IP = "127.0.0.1"
	pbb, _ := proto.Marshal(reg)
	println("pb序列化之后的长度：", len(pbb))

	outreg := &protoc.Register{}
	proto.Unmarshal(pbb, outreg)
	println(outreg.String())

	fmt.Println("pb=", reg.String())
	if outreg.SN == reg.SN && outreg.CompanyID == reg.CompanyID && outreg.ShopCode == reg.ShopCode {
		t.Log("TestPB测试通过")
	} else {
		t.Error("TestPB测试失败，因为输入与输出不符。")
	}
}

func TestPBCompress(t *testing.T) {
	reg := &protoc.Register{}
	reg.SN = "SN00112343"
	reg.CompanyID = 12345674879
	reg.ShopCode = "这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的这是测试之前的"
	reg.IP = "127.0.0.1"
	pbb, _ := proto.Marshal(reg)
	println("pb序列化之后的长度：", len(pbb))

	outreg := &protoc.Register{}
	proto.Unmarshal(pbb, outreg)
	println(outreg.String())

	fmt.Println("pb=", reg.String())
	if outreg.SN == reg.SN && outreg.CompanyID == reg.CompanyID && outreg.ShopCode == reg.ShopCode {
		t.Log("TestPBCompress测试通过")
	} else {
		t.Error("TestPBCompress测试失败，因为输入与输出不符。")
	}

}
