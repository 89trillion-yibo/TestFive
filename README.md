# 技术文档

### 1.整体框架

基于websocket协议开发一个在线聊天服务，通信使用protobuf进行通行，使用chan通道协助实现，使单个用户发消息全体在线用户可以接收，维护一个用户在线列表显示在线用户，并且支持用户登录和断开连接时全局消息提醒，支持心跳检测机制



### 2.目录结构

```
├── app
│   ├── httpserver
│   │   └── wsService.go
│   └── mian.go
├── go.mod
├── go.sum
├── internal
│   ├── ctrl
│   │   └── upgrade.go
│   ├── route
│   │   └── router.go
│   └── ws
│       ├── client.go
│       └── server.go
├── locust.py
├── logfile
│   ├── log.go
│   └── logInfo.txt
├── protobuf
│   ├── message.pb.go
│   └── message.proto
└── report.html

```



### 3.代码逻辑分层

| 层      | 文件夹                      | 主要职责                        | 调用关系                  | 其它说明     |
| ------- | --------------------------- | ------------------------------- | ------------------------- | ------------ |
| 应用层  | app/httpserver/wsService.go | 启动服务器                      | 调用路由层                | 不可同层调用 |
| 路由层  | internal/route/router.go    | 转发路由                        | 被应用层调用，调用控制层  | 不可同层调用 |
| 控制层  | internal/ctrl/upgrade.go    | 升级http请求为web socket        | 被路由层调用，调用service | 不可同层调用 |
| message | /protobuf                   | 存放proto相关文件               | 被handler、ctrl调用       | 不可同层调用 |
| 工具层  | /logfile                    | 存放日志相关                    | 被其他层调用              | 不可同层调用 |
| model   | /model                      | 存放数据结构与错误码            | 被其他层调用              | 不可同层调用 |
| WS      | internal/ws                 | 处理websocket连接的读写管理连接 | 被ctrl调用                | 可用层调用   |



### 4.存储设计

消息数据

| 内容         | 数据类形 | Key         |
| ------------ | -------- | ----------- |
| 消息内容     | string   | MessageText |
| 消息类型     | string   | MessageType |
| 用户名       | string   | User        |
| 在线用户列表 | []string | UserList    |



### 5.接口设计

##### 请求方式

websocket协议请求方式

##### 请求接口

ws://localhost:8080/ws

##### 请求参数

消息数据类型的二进制形式

##### 响应参数

二进制形式数据



### 6.第三方库

### gin

```
https://github.com/gin-gonic/gin
```

### Gorilla WebSocket

```
http://github.com/gorilla/websocket
```

### protobuf

```
http://github.com/golang/protobuf/proto
```



### 7.如何编译执行

```
//编译成可执行文件
 go build ./app/mian.go
 //运行文件
 ./mian
```



### 8.流程图

![未命名文件 (5)](https://user-images.githubusercontent.com/87186547/126989118-4b0f9924-c837-466f-bb88-f5d1efe1505b.jpg)
