package ws

import (
	"awesomeProject/Testfive/logfile"
	"awesomeProject/Testfive/protobuf"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
)

var (
	talk = "talk"
	exit = "exit"
	userlist = "userlist"
)

type ClientManagement struct {
	UserList map[*Client]bool //map存在线用户
	Talk     chan []byte      //广播发送消息
	NewSign  chan *Client     //新用户登陆
	Exit     chan *Client     //用户退出
	Single   chan *Client     //单播
}

func NewClientManagement() *ClientManagement{
	return &ClientManagement{
		UserList: make(map[*Client]bool),
		Talk:     make(chan []byte),
		NewSign:  make(chan *Client),
		Exit:     make(chan *Client),
		Single:   make(chan *Client),
	}
}


func (cl *ClientManagement) StartRun()  {
	fmt.Println("启动成功，开启监听")
	for {
		select {
		//新用户连接
		case conn := <- cl.NewSign:
			cl.UserList[conn] = true
			message := protobuf.Message{
				MessageText: "New user connection is " + conn.Name,
				MessageType: talk,
			}
			logfile.Info.Println("新的用户连接,"+conn.Name)
			marshal, _ := proto.Marshal(&message)
			cl.send(marshal)
		//用户退出
		case conn := <- cl.Exit:
			if _,ok :=cl.UserList[conn];ok{
				close(conn.SendMessage)
				delete(cl.UserList,conn)
				message := protobuf.Message{
					MessageText: conn.Name + " is Exit",
					MessageType: talk,
				}
				logfile.Info.Println(conn.Name,"用户退出")
				marshal, _ := proto.Marshal(&message)
				cl.send(marshal)
			}
		//返回用户列表
		case conn := <- cl.Single:
			userNameList := []string{}
			for c := range cl.UserList {
				userNameList = append(userNameList,c.Name)
			}
			message := protobuf.Message{
				MessageType: userlist,
				UserList: userNameList,
			}
			logfile.Info.Println("返回用户列表:",message)
			marshal, _ := proto.Marshal(&message)
			conn.SendMessage <- marshal
		//广播消息
		case talkMessage := <- cl.Talk:
			var meg string
			json.Unmarshal(talkMessage,&meg)
			message := protobuf.Message{
				MessageText: meg,
				MessageType: talk,
			}
			marshal, _ := proto.Marshal(&message)
			for c := range cl.UserList {
				select {
				case c.SendMessage <- marshal:
					logfile.Info.Println("发送消息:",meg)
					fmt.Println("message:",message)
				default:
					close(c.SendMessage)
					delete(cl.UserList,c)
					logfile.Error.Println("发消息错误，关闭连接")
				}
			}
		}
	}
}


//定义客户端管理的send方法
func (manager *ClientManagement) send(message []byte) {
	for conn := range manager.UserList {
			conn.SendMessage <- message
	}
}

