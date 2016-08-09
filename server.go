package main

import (
	"encoding/json"
	"tcp_long_connection/ConnMng"
	. "tcp_long_connection/config"
	"tcp_long_connection/logger"
	"tcp_long_connection/logic"
	"tcp_long_connection/tcpserver"
)

//tcp服务运行函数
func TcpRun() {
	server := tcpserver.New(Server.TcpIp + ":" + Server.TcpPort)
	server.OnNewClient(newClientConnect)
	server.OnNewMessage(newMessageComing)
	server.OnClientConnectionClosed(clientConnectClosed)
	//go detectEmptyConnection()
	server.Listen()
}

//新连接建立回调函数
func newClientConnect(c *tcpserver.Client) {
	// new client connected
	clientInfo := c.GetClientAddr()
	logger.Info("new client : %s\n", clientInfo)
	ConnMng.AddConnection(c)
}

//收到客户端消息回调函数
func newMessageComing(c *tcpserver.Client, p tcpserver.Package, message string) {
	// new message receive
	logger.Info("length := ", p.Length, " version := ", p.Version, " flag := ", p.Flag, " serial:= ", p.Serial)
	logger.Info(message)
	request := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &request)
	if err != nil {
		logger.Error(err.Error())
		c.Send(p, "{\"status\":1,\"message\":\"param error,json format error\"}\n")
		return
	}
	if request["cmd"] == nil || request["data"] == nil {
		c.Send(p, "{\"status\":1,\"message\":\"param error,no cmd or data\"}\n")
		return
	}
	cmd := ""
	ok := false
	if cmd, ok = request["cmd"].(string); !ok {
		c.Send(p, "{\"status\":1,\"message\":\"param error,cmd not string\"}\n")
		return
	}
	dataMap := make(map[string]interface{})
	if dataMap, ok = request["data"].(map[string]interface{}); !ok {
		c.Send(p, "{\"status\":1,\"message\":\"param error,data not map\"}\n")
		return
	}
	res, err := logic.Router(cmd, dataMap)
	if err != nil {
		c.Send(p, "{\"status\":1,\"message\":\""+err.Error()+"\"}\n")
		return
	}
	err = c.Send(p, res+"\n")
	if err != nil {
		logger.Error(err.Error())
	}
}

//客户端连接关闭回调函数
func clientConnectClosed(c *tcpserver.Client, err error) {
	// connection with client lost
	logger.Error("client closed" + c.GetClientAddr())
	ConnMng.RemoveConnection(c)
}
