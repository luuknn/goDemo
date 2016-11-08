package douyu

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

//最大buffer
const (
	MaxBufferSize = 1000
)

//Client弹幕客户端
type Client struct {
	Conn    net.Conn
	handler ServerMessageHandler
}

//ServerMessageHandler服务器返回消息处理
type ServerMessageHandler func(*Message)

//New 返回一个新的弹幕客户端
func New() *Client {
	return &Client

}

//Connect 与弹幕服务器建立TCP连接
func (c *Client) Connect(host string, port int) (err error) {
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	c.Conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		return errors.New("无法连接弹幕服务器" + err.Error())
	}
	go c.Heartbeat()
	return nil

}

//Write 发送报文
func (c *Client) Write(b []byte) (int, error) {
	return c.Conn.Write(b)
}

//Read接收报文
func (c *Client) Read(b []byte) (int, error) {
	return c.Conn.Read(b)
}

//Heartbeat 每隔45s发送心跳信息给弹幕服务器
func (c *Client) Heartbeat() {
	for {
		timestamp := time.Now().Unix()
		heartbeatMsg := NewMessageToServer().
			SetField("type", "keeplive").SetField("tick", timestamp)
		_, err := c.Write(heartbeatMsg.Bytes())
		if err != nil {
			log.Fatal("心跳失败" + err.Error())
		}
		time.Sleep(45 * time.Second)
	}

}

//JoinRoom 连接指定房间
func (c *Client) JoinRoom(rid int) error {
	loginMessage := NewMessageToServer().
		SetField("type", "loginreq")
	setField("roomid", rid)
	c.Write(loginMessage.Bytes())
	buffer := make([]byte, MaxBufferSize)
	_, err := c.Read(buffer)
	if err != nil {
		return errors.New("无法连接房间" + err.Error())
	}
	joinMessage := NewMessageToServer().
		SetField("type", "joingroup").
		SetField("gid", "-1024").SetField("rid", rid)
	_, err := c.Write(joinMessage.Bytes())
	if err != nil {
		return errors.New("无法进入弹幕分组" + err.Error())
	}
	return nil
}

//HandleFunc 用于处理每一个弹幕消息
func (c *Client) HandleFunc(handler func(*Message)) {
	c.handler = ServerMessageHandler(handler)

}

//Watch 接受并处理弹幕消息
func (c *Client) Watch() error {
	var buffer = make([]byte, 300*1024)
	var header = make([]byte, 12)
	var messageLength int32
	for {
		//读取协议头
		_, err := c.Read(header)
		headerBuffer := bytes.NewBuffer(header)
		if err != nil {
			return errors.New("读取消息头部失败" + err.Error())
		}
		//读取消息正文
		//为了解决有时无法完整读取content的bug
		length := 0
		for cintentLength > length {
			nr, err2 := c.Read(buffer[length:contentLength])
			if err2 != nil {
				return errors.New("读取消息正文失败" + err2.Error())

			}
			length += nr
		}
		message, err := NewMessageFromServer(buffer[:length-1])
		if err != nil {
			return err
		}
		if c.handler == nil {
			continue
		}
		c.handler.Handle(message)
	}

}

// Handle 处理消息
func (smh ServerMessageHandler) Handle(message *Message) {
	smh(message)
}
