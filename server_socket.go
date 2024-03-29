/*
 * @Author: F1
 * @Date: 2020-07-14 21:16:18
 * @LastEditors: F1
 * @LastEditTime: 2021-09-07 21:24:07
 * @Description:
 *
 *				yoyoecs　主要应用场景是边缘端与云端通讯时，采用socket来同步数据，该项目主要为底层协议及通讯实现。应最大限度的避开业务逻辑。
 *				核心为三大部分:
 *					第一部份为协议：ＰＲＯＴＯＣＯＬＳ中对头部、指令、标识位的定义
 *
 *						Header,Command,Flag
 *
 *					第二部份是客户端：
 *						ClientSocket 客户端对象，含连接及连接对象的状态，实现了重连机制，将收、发消息通过事件通知业务端。
 *
 *					第三部份是服务端：
 *						ServerSocket　服务端监听对象，含客户端连接的管理
 *
 */
package yoyoecs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs/protoc"
	"github.com/fengyily/yoyoecs/protocols"
	"github.com/fengyily/yoyoecs/utils"
	"github.com/golang/protobuf/proto"
)

var syncMessage map[string]chan []byte

func init() {
	syncMessage = make(map[string]chan []byte)
}

/**
 * @Title: ServerSocket
 * @Description:
 *
 * 				服务端监听对象，含客户端连接的管理
 *				应用场景为，云端启动时，采用该对象对指定端口进行监听，等待客户端来连接
 *
 * @Author: F1
 * @Date: 2020-07-21 11:18:46
 */
type ServerSocket struct {
	shutdown bool
	conn     *net.TCPListener

	OnConnect       func(string, *ClientSocket)
	OnRecvMessage   func(protocols.Header, []byte, *ClientSocket)
	OnSendToMessage func(string, protocols.Header, []byte, *ClientSocket)
	OnClose         func(string)
	OnError         func(*ClientSocket)
	OnRecvError     func(error)
	OnSendError     func(error)

	// 广播消息时的队列
	DataChan chan []byte

	Clients   map[string]*ClientSocket
	cloneLock sync.Mutex
}

/**
 * @Title:Run
 * @Description:
 *
 *						开始监听端口，等待边缘端连接
 *
 *						examples:
 *							- Run("*:9091")　开始监听端口，等待边缘端连接
 *
 * @Author: F1
 * @Date: 2020-07-21 11:41:22
 * @Param: address string
 * @Return: ok bool, err error
 */
func (server *ServerSocket) Run(address string) (ok bool, err error) {
	server.DataChan = make(chan []byte, 1000)
	server.Clients = make(map[string]*ClientSocket)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return false, err
	}

	server.conn, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return false, err
	}
	go server.send()

	go func() {
		for {
			if server.shutdown {
				fmt.Println("ServerSocket::Run", "服务已挂，通知Ｆ１处理吧")
				break
			}
			client, err := server.conn.Accept()
			if err != nil {
				fmt.Println("error", err)
				continue
			}
			cip := client.RemoteAddr().String()
			fmt.Println(cip, "连接成功。")

			c := &ClientSocket{
				RemoteAddr: cip,
				OnRecvMessage: func(h protocols.Header, d []byte, cs *ClientSocket) {
					if h.Cmd != protocols.REQUEST_HEARTBEAT {
						if _, ok := syncMessage[cs.ConnectId]; ok {
							syncMessage[cs.ConnectId] <- d
							return
						}
					}
					server.OnRecvMessage(h, d, cs)
				},
				OnSendToMessage: func(sn string, sendTo *protoc.SendTo, cs *ClientSocket) {
					fmt.Println(" server 收到了消息转发 SN:", cs.ConnectId, "send to :", sendTo.CID, "seq:", sendTo.Seq, "长度：", len(sendTo.Data))
					server.SendToByClientId(sendTo.CID, sendTo.Data)
					syncMessage[sendTo.CID] = make(chan []byte)
					go func() {
						if sendTo.Timeout > 60 {
							sendTo.Timeout = 60
						}
						ctx, _ := context.WithTimeout(context.TODO(), time.Duration(sendTo.Timeout)*time.Second)

						reply := protoc.SendToReply{
							Seq: sendTo.Seq,
						}
						fmt.Println(fmt.Sprintf("同步消息，在此等待:%s 回应", sendTo.CID))
						select {
						case reply.Data = <-syncMessage[sendTo.CID]:
							d, _ := proto.Marshal(&reply)
							replyData, _ := cs.InitMessage(protocols.REQUEST_SENDTO_REPLY, protocols.HEADER_FLAG_DATA_TYPE_PB, d)
							fmt.Println("收到同步返回的消息，长度为：", len(replyData))
							cs.SendData(replyData)
						case <-ctx.Done():
							fmt.Println("超时退出")
							cs.SendMessage(protocols.TIMEOUT, protocols.HEADER_FLAG_DATA_TYPE_PB, nil)
						}

						fmt.Println("任务完成，将同步消息删除")
						defer delete(syncMessage, sendTo.CID)
					}()
				},
				OnError: func(cs *ClientSocket) {
					if server.OnError != nil {
						server.OnError(cs)
					}
					fmt.Println("OnError")
				},
			}
			c.FormConn(&client)

			if server.OnConnect != nil {
				server.OnConnect(client.RemoteAddr().String(), c)
			}
		}
		fmt.Println("退出。。。")
	}()

	return true, nil
}

/**
 * @Title: 　RemoveClient
 * @Description:
 *
 *							将端从列表中移除
 *
 * @Author: F1
 * @Date: 2020-07-21 11:42:47
 * @Param:clientId string
 */
func (server *ServerSocket) RemoveClient(clientId string) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()

	cs, ok := server.Clients[clientId]
	if ok {
		delete(server.Clients, cs.ConnectId)
		fmt.Println("删除了连接：", cs.ConnectId)
		cs = nil
	} else {
		fmt.Println("来晚一步，早已被删除，当然，也可能从来就不曾有过")
	}
}

/**
 * @Title: ClientOnline
 * @Description:
 *
 *				客户端的在线状态，一般通过心跳来更新
 *
 * @Author: F1
 * @Date: 2020-07-21 11:42:47
 * @Param:clientId string, cs *ClientSocket
 */
func (server *ServerSocket) ClientOnline(clientId string, cs *ClientSocket) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()

	cs.OnlineTime = time.Now()

	if server.Clients != nil {
		server.Clients[clientId] = cs
	}
}

/**
 * @Title: AddClient
 * @Description:
 *
 *				将端添加到列表中
 *
 * @Author: F1
 * @Date: 2020-07-21 11:42:47
 * @Param:clientId string, cs *ClientSocket
 */
func (server *ServerSocket) AddClient(clientId string, cs *ClientSocket) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()
	if server.Clients != nil {
		cs.ConnectId = clientId
		server.Clients[clientId] = cs
	}
}

//
/**
 * @Title: SendMessage
 * @Description:
 *
 *				给所有的端点发送消息
 *
 * @Author: F1
 * @Date: 2020-07-21 11:42:47
 * @Param:header protocols.Header, data []byte
 */
func (server *ServerSocket) SendMessage(header protocols.Header, body []byte) {
	// 发送前对头部标识进行处理：压缩
	if (protocols.HEADER_FLAG_IS_COMPRESS&header.Flag) > 0 && len(body) > 0 {
		fmt.Println("发送消息：开启了压缩,压缩前", len(body))
		body = utils.Compress(body)
		fmt.Println("发送消息：开启了压缩,压缩后", len(body))
	}

	var data []byte
	if body != nil {
		header.Length = uint32(len(body))
		data = header.ToBytes()
		data = append(data, body...)
		fmt.Println("SendMessage cmd", header.Cmd, "length", header.Length, len(data))
	} else {
		data = header.ToBytes()
		fmt.Println("SendMessage cmd", header.Cmd, "length nil", header.Length, len(data))
	}

	//丢到队列中云处理
	server.DataChan <- data
}

/**
 * @Title: send
 * @Description:
 *
 *				发送，处理等待发送状态
 *
 * @Author: F1
 * @Date: 2020-07-21 11:44:57
 */
func (server *ServerSocket) send() {
	for {
		data := <-server.DataChan
		if server.shutdown {
			fmt.Println("服务已挂，可能是重启了，如果不是预期的结果，通知Ｆ１协助处理吧")
			break
		}
		go func() {
			clis := server.cloneTags(server.Clients)
			for _, c := range clis {
				fmt.Println("开始发送数据")
				go c.SendData(data)
			}
		}()
	}
}

/**
 * @Title: SendByClientId
 * @Description:
 *
 *				发送给某一个边缘端
 *
 * @Author: F1
 * @Date: 2020-07-21 11:45:33
 * @Param: clientId string, cmd protocols.Command, flag uint8, data []byte
 * @Return: err error)
 */
func (server *ServerSocket) SendByClientId(clientId string, cmd protocols.Command, flag protocols.Flag, data []byte) (err error) {
	client, ok := server.Clients[clientId]
	if ok {
		err = client.SendMessage(cmd, flag, data)
		fmt.Println("SendByClientId client.SendMessage(cmd, flag, data)", err)
	} else {
		err = errors.New("连接不存在啊，你确定它的状态是对的吗？")
	}

	return
}

func (server *ServerSocket) SendToByClientId(clientId string, data []byte) (err error) {
	client, ok := server.Clients[clientId]
	if ok {
		err = client.SendData(data)
		fmt.Println("SendToByClientId client.SendData(data)", err)
	} else {
		err = errors.New("连接不存在啊，你确定它的状态是对的吗？")
	}

	return
}

func (server *ServerSocket) cloneTags(tags map[string]*ClientSocket) map[string]*ClientSocket {
	cloneTags := make(map[string]*ClientSocket)
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()
	for k, v := range tags {
		cloneTags[k] = v
	}
	return cloneTags
}

/**
* @Title: Close
* @Description:
*
*				关闭监听对象，退出
*
* @Author: F1
* @Date: 2020-07-21 12:32:55
 * @Return: err error
*/
func (server *ServerSocket) Close() (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Close")
		}
	}()
	server.shutdown = true
	if server.conn != nil {
		err = server.conn.Close()
		fmt.Println(err)
	} else {
		fmt.Println(server.conn)
	}
	for _, v := range server.Clients {
		(*v.conn).Close()
	}
	server.conn = nil
	if err != nil {
		fmt.Println("连接关闭失败。", err)
	}
	fmt.Println("连接关闭成功。")
	server.Clients = nil
	close(server.DataChan)
	return
}
