package yoyoecs

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"github.com/fengyily/yoyoecs/protocols"
)

type ServerSocket struct {
	conn      net.Conn
	ipAddress string

	OnConnect     func(string)
	OnRecvMessage func(protocols.Header, []byte, *ClientSocket)
	OnClose       func(string)
	OnError       func(*ClientSocket)
	OnRecvError   func(error)
	OnSendError   func(error)

	DataChan chan []byte
	Buffer   []byte

	Clients   map[string]*ClientSocket
	cloneLock sync.Mutex
}

// Run("*:9091")
func (server *ServerSocket) Run(address string) (ok bool, err error) {
	server.DataChan = make(chan []byte, 1000)
	server.Clients = make(map[string]*ClientSocket)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return false, err
	}

	conn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return false, err
	}
	go server.send()
	go func() {
		for {
			client, err := conn.Accept()
			if err != nil {
				fmt.Println("error", err)
				continue
			}
			fmt.Println(client.RemoteAddr().String(), "连接成功。")

			c := &ClientSocket{
				OnRecvMessage: func(header protocols.Header, data []byte, cs *ClientSocket) {
					if server.OnRecvMessage != nil {
						server.OnRecvMessage(header, data, cs)
					}
				},
				OnError: func(cs *ClientSocket) {
					if server.OnError != nil {
						server.OnError(cs)
					}
				},
			}
			c.FormConn(&client)

			if server.OnConnect != nil {
				server.OnConnect(client.RemoteAddr().String())
			}
		}
	}()

	return true, nil
}

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

func (server *ServerSocket) AddClient(clientId string, cs *ClientSocket) {
	server.cloneLock.Lock()
	defer server.cloneLock.Unlock()

	server.Clients[clientId] = cs
}

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
		go func() {
			clis := server.cloneTags(server.Clients)
			for _, c := range clis {
				fmt.Println("开始发送数据")
				c.SendData(data)
			}
		}()
	}
}

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
