/*
 * @Author: F1
 * @Date: 2020-07-21 11:47:32
 * @LastEditors: F1
 * @LastEditTime: 2020-07-30 15:22:12
 * @Description: 客户端测试
 */
package main

import (
	"fmt"
	"net"
	"strconv"
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
		fmt.Println("成功收到消息：", header.Cmd, "长度：", header.Length, "Flag", fmt.Sprintf("%08b", header.Flag), string(data))
		fmt.Println("-----------------------------------------------------------", time.Now())

		if header.Cmd == protocols.RESPONSE_UPLOAD_SKU_DATA {
			fmt.Println("通知继续发送。")
		}
		if header.Cmd == protocols.RESPONSE_REGISTER_SUCCESS {
			fmt.Println("收到注册成功消息")

			//test(cs)
		}
	}

	client.OnConnect = func(ip string, cs *yoyoecs.ClientSocket) {
		// 采用ＰＢ上传ＳＫＵ信息的示例
		info := protoc.Register{}
		info.CompanyID = 123456789
		info.ShopCode = "shopcode123456789测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈"
		d, _ := proto.Marshal(&info)
		client.SendMessage(protocols.REQUEST_REGISTER, protocols.HEADER_FLAG_IS_COMPRESS|protocols.HEADER_FLAG_DATA_TYPE_JSON, d)
		fmt.Println("发起了注册申请")

		skulist := protoc.SkuList{
			Sku: make([]*protoc.Sku, 0),
		}

		for i := 0; i < 200; i++ {
			sku := &protoc.Sku{}
			sku.Id = int32(i)
			sku.SkuName = "skuname" + strconv.Itoa(i)
			sku.Price = 100

			skulist.Sku = append(skulist.Sku, sku)
		}

		body, _ := proto.Marshal(&skulist)

		//　采用了pb数据类型，并且开启压缩，如果传输数据较大，建议开启压缩
		client.SendMessage(protocols.REQUEST_UPLOAD_SKU_DATA, protocols.HEADER_FLAG_DATA_TYPE_PB|protocols.HEADER_FLAG_IS_COMPRESS, body)

		item := protoc.ItemList{
			Items: make([]*protoc.Item, 0),
		}

		for i := 0; i < 100; i++ {
			it := protoc.Item{}
			it.Id = int64(i)
			it.Name = "测试啊"
			it.MatchVersionCode = "v1.0.0.1"

			item.Items = append(item.Items, &it)
		}
		body, _ = proto.Marshal(&item)
		client.SendMessage(protocols.REQUEST_TRANS_ITEM_DATA, protocols.HEADER_FLAG_DATA_TYPE_PB|protocols.HEADER_FLAG_IS_COMPRESS, body)

		yoyoList := protoc.YoyoInfoList{
			YoyoInfo: make([]*protoc.YoyoInfo, 0),
		}
		for i := 0; i < 100; i++ {
			yoyoInfo := protoc.YoyoInfo{}
			yoyoInfo.Name = "Name" + strconv.Itoa(i)

			yoyoList.YoyoInfo = append(yoyoList.YoyoInfo, &yoyoInfo)
		}
		body, _ = proto.Marshal(&yoyoList)
		client.SendMessage(protocols.REQUEST_TRANS_YOYOINFO_DATA, protocols.HEADER_FLAG_DATA_TYPE_PB|protocols.HEADER_FLAG_IS_COMPRESS, body)

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

	go func(client *yoyoecs.ClientSocket) {

		index := 1
		n := 0
		for i := 1; i <= 100; i++ {
			for j := 1; j <= 100; j++ {

				index++
				client.SendMessage(protocols.RESPONSE_TRANS_SKU_DATA, protocols.HEADER_FLAG_IS_COMPRESS, []byte("「"+strconv.Itoa(index)+"」这是ＳＫＵ信息，shopcode123456789测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈测试的哈第"+strconv.Itoa(i)+" 轮测试"))
				if n%20 == 0 {
					index++
					SendSplitMessage(client.GetConn(), protocols.REQUEST_UPLOAD_SKU_DATA, []byte("「"+strconv.Itoa(index)+"」这是ＳＫＵ信息，第"+strconv.Itoa(i)+" 批分裂消息"))
				}
				if n%100 == 0 {
					index++
					SendBadMessage(client.GetConn(), protocols.REQUEST_UPLOAD_SKU_DATA, []byte("「"+strconv.Itoa(index)+"」这是一条坏消息，它将影响下一条消息。就看你的程序会不会识别，不要影响下下条消息。。。。"))
				}
				n++
			}
			time.Sleep(2 * time.Second)
		}
	}(client)

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
