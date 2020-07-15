package yoyoecs

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fengyily/yoyoecs/protocols"
)

//　ServerSocket　服务端连接对象
type ServerSocket struct {
	shutdown bool
	conn *net.TCPListener

	OnConnect     func(string, *ClientSocket)
	OnRecvMessage func(protocols.Header, []byte, *ClientSocket)
	OnClose       func(string)
	OnError       func(*ClientSocket)
	OnRecvError   func(error)
	OnSendError   func(error)

	// 广播消息时的队列
	DataChan chan []byte

	Clients   map[string]*ClientSocket
	cloneLock sync.Mutex
}

// Run("*:9091")　开始监听端口，等待边缘端连接
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
				break
			}
			client, err := server.conn.Accept()
			if err != nil {
				fmt.Println("error", err)
				continue
			}
			fmt.Println(client.RemoteAddr().String(), "连接成功。")

			c := &ClientSocket{
				OnRecvMessage: server.OnRecvMessage,
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

//　RemoveClient　将端从列表中移除
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

func (server *ServerSocket) ClientOnline(clientId string, cs *ClientSocket) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()

	cs.OnlineTime = time.Now()
	
	if server.Clients != nil {
		server.Clients[clientId] = cs
	}
}

// AddClient 将端添加到列表中
func (server *ServerSocket) AddClient(clientId string, cs *ClientSocket) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()
	if server.Clients != nil {
		server.Clients[clientId] = cs
	}
}

// SendMessage　给所有的端点发送消息
func (server *ServerSocket) SendMessage(header protocols.Header, data []byte) {
	if data != nil {
		header.Length = uint16(len(data))
	}
	h := header.ToBytes()
	if header.Length > 0 {
		d := append(h, data...)
		server.DataChan <- d
	} else {
		server.DataChan <- h
	}
}

func (server *ServerSocket) send() {
	for {
		data := <-server.DataChan
		if server.shutdown {
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

// SendByClientId 发送给某一个边缘端
func (server *ServerSocket) SendByClientId(clientId string, cmd protocols.Command, flag uint8, data []byte) (err error) {
	client, ok := server.Clients[clientId]
	if ok {
		header := protocols.Header{}
		header.Cmd = cmd
		header.Flag = flag
		header.Length = uint16(len(data))

		h := header.ToBytes()
		if header.Length > 0 {
			d := append(h, data...)
			err = client.SendData(d)
		} else {
			err = client.SendData(h)
		}
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
		v.conn.Close()
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
