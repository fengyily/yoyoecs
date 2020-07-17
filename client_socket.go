package yoyoecs

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs/protocols"
	"github.com/fengyily/yoyoecs/utils"
)

type ClientSocket struct {
	isReConnect   bool
	ConnectId     string
	IsConnected   bool
	ipAddress     string
	conn          net.Conn
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

func (cs *ClientSocket) GetConn() (conn *net.Conn) {
	return &cs.conn
}

func (cs *ClientSocket) FormConn(conn *net.Conn) {
	cs.IsConnected = true
	cs.isReConnect = false
	cs.conn = *conn
	go cs.read()
}

func (cs *ClientSocket) Conn(ipAddress string) (err error) {
	cs.DataChan = make(chan []byte, 1000)
	cs.Buffer = make([]byte, 0)
	cs.IsConnected = false
	cs.isReConnect = true
	for {
		cs.ipAddress = ipAddress
		cs.conn, err = net.Dial("tcp", ipAddress)
		if err != nil {
			if cs.OnConnError != nil {
				cs.OnConnError(err)
			}
			time.Sleep(5 * time.Second)
			continue
		} else {
			fmt.Println("success")
			if cs.OnConnect != nil {
				cs.OnConnect("连接成功。", cs)
			}

			//cs.Register()
			cs.IsConnected = true
			break
		}
	}
	cs.RunHeartBeat.Do(cs.HeartBeat)

	go cs.read()

	return
}

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
		cs.conn, err = net.Dial("tcp", cs.ipAddress)

		if err != nil {
			if cs.OnConnError != nil {
				cs.OnConnError(err)
			}
			time.Sleep(5 * time.Second)
			continue
		} else {
			cs.OnConnect("连接成功。", cs)
			cs.IsConnected = true
			break
		}
	}
	return
}

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

func (cs *ClientSocket) connerror(err error) {
	if cs.conn != nil {
		cs.conn.Close()
		fmt.Println("连接断开")
	}

	cs.IsConnected = false
	if cs.OnError != nil {
		cs.OnError(cs)
	}
}

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
				break
			}
		}
		if cs.conn != nil {
			n, err := cs.conn.Read(data)
			if err != nil {
				cs.connerror(err)
				if cs.isReConnect {
					continue
				} else {
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
					fmt.Println("收到心跳回复。")
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
							if protocols.Flag(header.Flag)&protocols.HEADER_FLAG_IS_COMPRESS > 0 {

								fmt.Println("开启了压缩,解压前", len(data))
								data = utils.UnCompress(data)
								fmt.Println("开启了压缩,解压后", len(data))
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

func (cs *ClientSocket) SendMessage(cmd protocols.Command, flag byte, body []byte) (err error) {
	header := protocols.Header{}
	header.Cmd = cmd
	header.Flag = flag
	if len(body) > 2<<15 {
		panic(fmt.Sprintf("超出可接收长度。len(body):%d > %d", len(body), 2<<15))
	}

	var data []byte
	if body != nil {
		if protocols.HEADER_FLAG_IS_COMPRESS&protocols.Flag(flag) > 0 {
			fmt.Println("开启了压缩,压缩前", len(body))
			body = utils.Compress(body)
			fmt.Println("开启了压缩,压缩后", len(body))
		}

		header.Length = uint16(len(body))
		data = header.ToBytes()
		data = append(data, body...)
		fmt.Println("SendMessage cmd", cmd, "length", header.Length, len(data))
	} else {
		data = header.ToBytes()
		fmt.Println("SendMessage cmd", cmd, "length nil", header.Length, len(data))
	}

	err = cs.SendData(data)
	return
}

func (cs *ClientSocket) SendData(body []byte) (err error) {
	if cs.conn == nil {
		cs.connerror(err)
		return
	}

	cs.sendLock.Lock()
	defer cs.sendLock.Unlock()

	total := len(body)
	index := 0
	for index < total {
		send, err := cs.conn.Write(body[index:])
		if err != nil {
			fmt.Println(err.Error())
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
	return nil
}
