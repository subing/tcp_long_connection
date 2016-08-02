package main

import (
	"fmt"
	"os"
	"strings"
	. "tcp_long_connection/config"
	"tcp_long_connection/logger"

	"github.com/xgdapg/daemon"
)

func main() {
	//config daemon and show usage
	if len(os.Args) >= 2 {
		if os.Args[1] == "daemon" {
			//后台启动及看门狗进程
			daemon.Exec(daemon.Daemon | daemon.Monitor)
		} else if os.Args[1] == "--help" {
			fmt.Printf("Usage : %s [--help | daemon]\n", os.Args[0])
			fmt.Printf("--help    Show this help info\n")
			fmt.Printf("daemon    Start Server as daemon\n")
			return
		}
	}

	//初始化配置信息
	InitServiceConfig()
	//设置日志前台可见
	consoleFlag := false
	//set log level and init logger
	if Log.IsConsole == "1" {
		consoleFlag = true
	}
	logger.SetConsole(consoleFlag)
	//根据配置文件设置日志等级
	logLevel := logger.ERROR
	if strings.EqualFold(Log.Level, "Debug") {
		logLevel = logger.DEBUG
	} else if strings.EqualFold(Log.Level, "Info") {
		logLevel = logger.INFO
	} else if strings.EqualFold(Log.Level, "Warn") {
		logLevel = logger.WARN
	} else if strings.EqualFold(Log.Level, "Error") {
		logLevel = logger.ERROR
	} else if strings.EqualFold(Log.Level, "Fatal") {
		logLevel = logger.FATAL
	}
	logger.SetLevel(logLevel)
	//根据配置文件，设置日志路径，日志名，日志切割大小限制
	logger.SetRollingFile(Log.Path, Log.FileName, int32(Log.FileCount), int64(Log.FileSize), logger.MB)
	TcpRun()
}
