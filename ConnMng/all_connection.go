package ConnMng

import (
	"errors"
	"sync"
	"tcp_long_connection/logger"
	"tcp_long_connection/tcpserver"
)

var listMutext = new(sync.Mutex)
var allConnMap = make(map[*tcpserver.Client]bool)

func AddConnection(client *tcpserver.Client) error {
	if client == nil {
		logger.Error("client is nil")
		return errors.New("client is nil")
	}
	listMutext.Lock()
	defer listMutext.Unlock()
	if _, ok := allConnMap[client]; !ok {
		allConnMap[client] = true
	}
	return nil
}

func RemoveConnection(client *tcpserver.Client) error {
	if client == nil {
		logger.Error("client is nil")
		return errors.New("client is nil")
	}
	listMutext.Lock()
	defer listMutext.Unlock()
	delete(allConnMap, client)
	return nil
}

func GetAllClients() map[*tcpserver.Client]bool {
	return allConnMap
}

func GetConnectionSize() int {
	return len(allConnMap)
}
