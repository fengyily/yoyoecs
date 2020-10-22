/*
 * @Author: F1
 * @Date: 2020-07-14 21:16:18
 * @LastEditors: F1
 * @LastEditTime: 2020-10-22 21:41:41
 * @Description:
 *
 *				yoyoecs　主要应用场景是边缘端与云端通讯时，采用socket来同步数据，该项目主要为底层协议及通讯实现。应最大限度的避开业务逻辑。
 *				核心为三大部分:
 *					第一部份为协议：protocols中对头部、指令、标识位的定义
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
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs/protocols"
	"github.com/fengyily/yoyoecs/utils"
)

/**
 * @Title: ClientSocket
 * @Description:
 *
 * 				客户端对象，含连接及连接对象的状态，实现了重连机制，将收、发消息通过事件通知业务端。
 *				应用场景为，边缘端连接云端服务时，采用该对象的Conn。云端服务接收到客户端请求时，采用FormConn初使化该对象与连接对应
 *
 * @Author: F1
 * @Date: 2020-07-21 11:18:46
 */
type ClientSocket struct {
	isReConnect   bool
	ConnectId     string
	IsConnected   bool
	ipAddress     string
	conn          *net.Conn
	OnConnect     func(string, *ClientSocket)
	OnRecvMessage func(protocols.Header, []byte, *ClientSocket)
	OnClose       func(string)
	OnError       func(*ClientSocket)
	OnConnError   func(error)
	OnRecvError   func(error)
	OnSendError   func(error)

	DataChan     chan []byte
	Buffer       []byte
	RunHeartBeat sync.Once
	IsConnecting bool
	OnlineTime   time.Time
	ReConnLock   sync.Mutex
	sendLock     sync.Mutex
}

/**
* @Title: GetConn
* @Description:
*
*				获取连接对象
*
* @Author: F1
* @Date: 2020-07-21 11:23:03
 * @Return:conn *net.Conn
*/
func (cs *ClientSocket) GetConn() (conn *net.Conn) {
	return cs.conn
}

/**
 * @Title:FormConn
 * @Description:
 *
 *				从连接对象初使化客户端实例，以便于通过客户端实例来实现基于统一协议统收发消息
 *
 * @Author: F1
 * @Date: 2020-07-21 11:24:09
 * @Param:conn *net.Conn
 */
func (cs *ClientSocket) FormConn(conn *net.Conn) {
	cs.IsConnected = true
	cs.isReConnect = false
	cs.conn = conn
	go cs.read()
}

/**
 * @Title:RemoteAddr
 * @Description:
 *
 *				获取远端的ＩＰ地址
 *
 * @Author: F1
 * @Date: 2020-10-22 21:41:09
 */
func (cs *ClientSocket) RemoteAddr() string {
	return (*(cs.conn)).RemoteAddr().String()
}

/**
 * @Title: Conn
 * @Description:
 *
 *				边接服务端
 *
 * @Author: F1
 * @Date: 2020-07-21 11:25:32
 * @Param: ipAddress string
 * @Return: err error
 */
func (cs *ClientSocket) Conn(ipAddress string) (err error) {
	cs.DataChan = make(chan []byte, 1000)
	cs.Buffer = make([]byte, 0)
	cs.IsConnected = false
	cs.isReConnect = true
	for {
		cs.ipAddress = ipAddress
		cn, err := net.Dial("tcp", ipAddress)
		cs.conn = &cn
		if err != nil {
			if cs.OnConnError != nil {
				cs.OnConnError(err)
			}
			time.Sleep(5 * time.Second)
			continue
		} else {
			fmt.Println("success")
			cs.IsConnected = true
			if cs.OnConnect != nil {
				cs.OnConnect("连接成功。", cs)
			}
			break
		}
	}
	cs.RunHeartBeat.Do(cs.HeartBeat)

	go cs.read()

	return
}

/**
 * @Title: checkConn
 * @Description:
 *
 *				检查连接的状态，如果连接断开，会尝试重新连接，直到重连成功。
 *
 * @Author: F1
 * @Date: 2020-07-21 11:26:31
 * @Return: err error
 */
func (cs *ClientSocket) checkConn() (err error) {
	if cs.IsConnected {
		return
	}

	cs.ReConnLock.Lock()
	defer cs.ReConnLock.Unlock()

	cs.DataChan = make(chan []byte, 1000)
	cs.Buffer = cs.Buffer[len(cs.Buffer):]
	for {
		cs.conn = nil
		cn, err := net.Dial("tcp", cs.ipAddress)
		cs.conn = &cn
		if err != nil {
			if cs.OnConnError != nil {
				cs.OnConnError(err)
			}
			time.Sleep(5 * time.Second)
			continue
		} else {
			cs.IsConnected = true
			cs.OnConnect("连接成功。", cs)
			break
		}
	}
	return
}

/**
 * @Title: HeartBeat
 * @Description:
 *
 *				连接心跳，在连接的情况下，每５秒发送一次心跳，目前暂时是固定的
 *
 * @Author: F1
 * @Date: 2020-07-21 11:27:44
 */
func (cs *ClientSocket) HeartBeat() {
	go func() {
		for {
			if !cs.IsConnected {
				fmt.Println("连接已断开，尝试重连中")
				if cs.isReConnect {
					cs.checkConn()
				} else {
					break
				}

			}
			cs.SendMessage(protocols.REQUEST_HEARTBEAT, 0, nil)
			time.Sleep(5 * time.Second)
		}
	}()
}

/**
 * @Title: connerror
 * @Description:
 *
 *				连接异常的处理
 *
 * @Author: F1
 * @Date: 2020-07-21 11:29:23
 * @Param:err error
 */
func (cs *ClientSocket) connerror(err error) {
	if cs.conn != nil {
		(*cs.conn).Close()
		fmt.Println("连接断开　关闭连接")
	}

	cs.IsConnected = false
	if cs.OnError != nil {
		fmt.Println("连接断开 通知OnError")
		cs.OnError(cs)
	}
}

/**
 * @Title: RemoteIpAddress
 * @Description:
 *
 *				RemoteIpAddress 远端IP
 *
 * @Author: F1
 * @Date: 2020-07-21 11:29:23
 * @Param:err error
 */
func (cs *ClientSocket) RemoteIpAddress() string {
	if cs.conn != nil && cs.IsConnected {
		return (*cs.conn).RemoteAddr().String()
	}
	return "Unknow"
}

/**
 * @Title: read
 * @Description:
 *
 *				监控连接对象，开启循环读连接对象中的字节流，并解析成标准包输出给业务
 *
 * @Author: F1
 * @Date: 2020-07-21 11:30:16
 * @Param:
 * @Return:
 */
func (cs *ClientSocket) read() {
	fmt.Println("进入循环读")
	data := make([]byte, 4096)
	for {
		if !cs.IsConnected {
			time.Sleep(5 * time.Second)
			fmt.Println("连接断开")
			if cs.isReConnect {
				continue
			} else {
				fmt.Println("read", "连接断开了？？？")
				break
			}
		}
		if cs.conn != nil {
			n, err := (*cs.conn).Read(data)
			if err != nil {
				cs.connerror(err)
				if cs.isReConnect {
					continue
				} else {
					fmt.Println("read", "２－　连接断开了？？？")
					break
				}
			}
			// merge buffer
			if n > 0 {
				// 加上缓冲区的数据
				// befor := len(cs.Buffer)
				cs.Buffer = append(cs.Buffer, data[:n]...)
				//fmt.Println("原缓存右有", befor, "这次又读了", n, "befor string=", string(cs.Buffer[:befor]), "merge string=", string(data[0:n]))
			}

			for i := 0; i < len(cs.Buffer); {
				var header protocols.Header
				ok, header := protocols.LoadHeader(&cs.Buffer)
				if !ok {
					fmt.Println("加载头部失败。。。")
					break
				} else {
					//fmt.Println("加载头部成功：", header.Length, header.Cmd)
				}

				if header.Cmd == protocols.REQUEST_HEARTBEAT {
					// 收到心跳包，马上回复一个
					go cs.SendMessage(protocols.RESPONSE_HEARTBEAT, 0, nil)
					//cs.Buffer = cs.Buffer[i+1:]
					if cs.OnRecvMessage != nil {
						cs.OnRecvMessage(header, nil, cs)
					}
					continue
				} else if header.Cmd == protocols.RESPONSE_HEARTBEAT {
					//fmt.Println("收到心跳回复。")
					continue
				} else {
					total := len(cs.Buffer)
					// if header.Length == 0 {
					// 	//fmt.Println("指令对应的长度不对，重新定位丢弃１", string(cs.Buffer))
					// 	// cs.Buffer = cs.Buffer[i+1:]
					// 	// fmt.Println("length", header.Length, "len(cs.Buffer)", len(cs.Buffer), "i", i, "cmd", header.Cmd.ToString())

					// 	continue
					// }

					// 判断剩余长度是否满足头部描述长度。
					balance := total - i - protocols.HEADER_LENGTH

					if balance >= int(header.Length) {
						//fmt.Println(cs.OnRecvMessage)
						//fmt.Println(cs.Buffer[i+protocols.HEADER_LENGTH:i+protocols.HEADER_LENGTH+int(header.Length)])

						// 收完后，直接处理缓存内容
						data := cs.Buffer[i+protocols.HEADER_LENGTH : i+protocols.HEADER_LENGTH+int(header.Length)]
						cs.Buffer = cs.Buffer[i+protocols.HEADER_LENGTH+int(header.Length):]
						if cs.OnRecvMessage != nil {
							if header.Flag&protocols.HEADER_FLAG_IS_COMPRESS > 0 {

								//fmt.Println("收到消息：开启了压缩,解压前", len(data))
								data = utils.UnCompress(data)
								//fmt.Println("收到消息：开启了压缩,解压后", len(data))
								header.Length = uint16(len(data))
							}
							cs.OnRecvMessage(header, data, cs)

						}

						continue
					} else {
						// 如果不够一个完整的包，存入缓冲区
						//fmt.Println("不够一个完整的包，存入缓冲区", "总共：", total, "内容是：", string(cs.Buffer[protocols.HEADER_LENGTH:total]), "剩余：", balance, "包体长度是：", header.Length, i)
						break
					}
				}
			}
		} else {
			fmt.Println("连接可能中断了，可能是服务端挂了，但生活还得继续！等待时机，以不变应万变！！")
			continue
		}
	}
	fmt.Println("退出了循环读")
}

/**
 * @Title: SendMessage
 * @Description:
 *
 *				发送消息，从连接中发送消息
 *
 * @Author: F1
 * @Date: 2020-07-21 11:31:57
 * @Param: cmd protocols.Command 指令
 * @Param: flag byte　标识
 * @Param: body []byte　包体
 * @Return:err error
 */
func (cs *ClientSocket) SendMessage(cmd protocols.Command, flag protocols.Flag, body []byte) (err error) {
	header := protocols.Header{}
	header.Cmd = cmd
	header.Flag = flag

	if (protocols.HEADER_FLAG_IS_COMPRESS&flag) > 0 && len(body) > 0 {
		fmt.Println("发送消息：开启了压缩,压缩前", len(body))
		body = utils.Compress(body)
		fmt.Println("发送消息：开启了压缩,压缩后", len(body))
	}

	if len(body) > 2<<15 {
		//panic(fmt.Sprintf("超出可接收长度。len(body):%d > %d", len(body), 2<<15))
		fmt.Println(fmt.Sprintf("超出可接收长度。len(body):%d > %d", len(body), 2<<15))
		return
	}

	var data []byte
	if body != nil {
		header.Length = uint16(len(body))
		data = header.ToBytes()
		data = append(data, body...)
		//fmt.Println("SendMessage cmd", cmd, "length", header.Length, len(data))
	} else {
		data = header.ToBytes()
		//fmt.Println("SendMessage cmd", cmd, "length nil", header.Length, len(data))
	}

	err = cs.SendData(data)
	return
}

/**
 * @Title: SendData
 * @Description:
 *
 *				发送数据
 *
 * @Author: F1
 * @Date: 2020-07-21 11:33:47
 * @Param:body []byte
 * @Return:err error
 */
func (cs *ClientSocket) SendData(body []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendData 出现异常了，这个是需要查下原因的。", err)
		}
	}()
	if cs.conn == nil {
		fmt.Println("SendData", "连接异常，通知主人")
		cs.connerror(err)
		return
	}
	//fmt.Println("SendData", "准备发送，获取待锁。")
	cs.sendLock.Lock()
	defer cs.sendLock.Unlock()
	//fmt.Println("SendData", "准备发送，获取待锁成功。")
	total := len(body)
	index := 0

	// 确保body中的数据全部发送完成。
	for index < total {
		if cs.conn == nil || !cs.IsConnected {
			fmt.Println("SendData", "连接异常，怎么肥４？")
			//cs.connerror(err)
			break
		}
		send, err := (*cs.conn).Write(body[index:])
		if err != nil {
			fmt.Println("SendData", "发送异常，这个问题是严重的,可能会导致连接断开。", err.Error())
			if cs.OnSendError != nil {
				cs.OnSendError(err)
			}
			if cs.OnError != nil {
				cs.OnError(cs)
			}
			break
		}
		index += send
	}

	//fmt.Println("SendData", "成功发送：", total)
	return nil
}
