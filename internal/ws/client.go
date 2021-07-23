package ws

import (
	"awesomeProject/Testfive/protobuf"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"time"
)

const (
	//允许向客户端写入消息的时间
	writeWait = 10 * time.Second

	//允许从客户端读取下一个 pong 消息的时间。
	pongWait = 60 * time.Second

	//在此期间向对等方发送 ping。必须小于 pongWait。
	pingPeriod = (pongWait * 9) / 10

	//对等方允许的最大消息大小。
	maxMessageSize = 512
)

type Client struct {
	Conn *websocket.Conn           //用户websocket连接
	Name string                    //用户名称
	SendMessage chan []byte        //发送消息
	ClManagement *ClientManagement //连接管理
}

//心跳检测机制借鉴github上开源代码
//github地址：https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
//当服务器接收到客户端发来的消息在此协程处理
func (c *Client) Read() {
	defer func() {
		c.ClManagement.Exit <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	//设置底层网络连接的读取截止时间。
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for  {
		_, message, err := c.Conn.ReadMessage()

		if err != nil{
			c.ClManagement.Exit <- c
			c.Conn.Close()
			break
		}
		meg := protobuf.Message{}
		proto.Unmarshal(message, &meg)
		//err = proto.Unmarshal(message, &meg)
		fmt.Println("meg:",meg)
		if err!=nil{
			fmt.Println(err)
		}
		switch meg.MessageType {
		case "talk":
			fmt.Println("Talk:",meg.MessageText)
			talkMessage := meg.User+":"+meg.MessageText
			data, _ := json.Marshal(talkMessage)
			//广播发送
			c.ClManagement.Talk <- data
		case "exit":
			c.ClManagement.Exit <- c
		case "userlist":
			c.ClManagement.Single <- c
		}
	}
}

//服务端需要返回数据给客户端
func (c *Client) Write()  {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message,ok := <- c.SendMessage:
			if !ok{
				c.Conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage,message)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			fmt.Println("心跳检测")
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

