package ConnMng

import (
	"errors"
	"sync"
	"tcp_long_connection/logger"
	"tcp_long_connection/tcpserver"
)

var mutex = new(sync.RWMutex)
var clientInfoMap = make(map[string]*tcpserver.Client)

func AddClient(session string, client *tcpserver.Client) error {
	mutex.Lock()
	defer mutex.Unlock()
	if clientInfoMap[session] != nil {
		logger.Error("client : ", session, " exist")
		return errors.New("client : " + session + " exist")
	}
	clientInfoMap[session] = client
	return nil
}

func GetClient(session string) *tcpserver.Client {
	mutex.RLock()
	defer mutex.RUnlock()
	if clientInfoMap[session] == nil {
		logger.Error("client : ", session, " not exist")
		return nil
	}
	return clientInfoMap[session]
}

func RemoveClient(session string) {
	mutex.Lock()
	defer mutex.Unlock()
	if clientInfoMap[session] == nil {
		logger.Error("client : ", session, " not exist")
		return
	}
	delete(clientInfoMap, session)
}

func UpdateClient(session string, client *tcpserver.Client) {
	mutex.Lock()
	defer mutex.Unlock()
	clientInfoMap[session] = client
	return
}

func GetClientTotal() int {
	return len(clientInfoMap)
}
