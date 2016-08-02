package config

import (
	"strconv"

	"github.com/widuu/goini"
)

type ServerInfo struct {
	Ip      string
	Port    string
	TcpIp   string
	TcpPort string
}

type LogInfo struct {
	Level     string
	IsConsole string
	FileCount int
	FileSize  int
	Path      string
	FileName  string
	SockUrl   string
}

var Server ServerInfo
var Log LogInfo

func InitServiceConfig() error {
	conf := goini.SetConfig("config/config.conf")
	//read server config info
	Server.Ip = conf.GetValue("ServerInfo", "ip")
	Server.Port = conf.GetValue("ServerInfo", "port")
	Server.TcpIp = conf.GetValue("ServerInfo", "tcpip")
	Server.TcpPort = conf.GetValue("ServerInfo", "tcpport")

	//read log config
	Log.Level = conf.GetValue("LogInfo", "level")
	Log.Path = conf.GetValue("LogInfo", "path")
	Log.IsConsole = conf.GetValue("LogInfo", "isConsole")
	count, err := strconv.Atoi(conf.GetValue("LogInfo", "filecount"))
	if err != nil {
		count = 10
	}
	Log.FileCount = count
	size, err := strconv.Atoi(conf.GetValue("LogInfo", "filesize"))
	if err != nil {
		size = 100
	}
	Log.FileSize = size
	Log.FileName = conf.GetValue("LogInfo", "filename")
	Log.SockUrl = conf.GetValue("LogInfo", "sockurl")

	return nil
}
