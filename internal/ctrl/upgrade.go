package ctrl

import (
	"awesomeProject/Testfive/internal/ws"
	"awesomeProject/Testfive/logfile"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeHttp(cm *ws.ClientManagement, w http.ResponseWriter,r *http.Request) {
	//升级http请求
	if websocket.IsWebSocketUpgrade(r){
		conn, err := upgrader.Upgrade(w, r, nil)
		if err!=nil{
			fmt.Println(err)
		}
		client := &ws.Client{
			Conn: conn,
			Name: r.Header.Get("name"),
			SendMessage: make(chan []byte),
			ClManagement: cm,
		}
		fmt.Println("name",client.Name)
		//新创建的客户端连接进入通道
		client.ClManagement.NewSign <- client
		//开启服务端从通道读数据
		go client.Read()
		//开启服务端往通道写数据
		go client.Write()
	}else {
		logfile.Error.Println("发生错误")
		fmt.Println("错误")
	}
}