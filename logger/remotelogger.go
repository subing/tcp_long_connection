package logger

import (
	conf "KernelService/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//http log requester
func httpLog(level int, logid, vplNumber, data string) {
	logType := "exception"
	if level == INFO {
		logType = "info"
	}
	dataMap := map[string]interface{}{
		"logid":  logid,
		"vpln":   vplNumber,
		"source": "ApiServer",
		"type":   logType,
		"data":   data,
		"time":   time.Now().Unix()}
	logByteArry, err := json.Marshal(dataMap)
	if err != nil {
		Error(err.Error())
		return
	}
	logRes, err := SimplePost(conf.Log.SockUrl, string(logByteArry))
	if err != nil {
		return
	}
	Debug("HttpLog Response : ", logRes)
}

func SimplePost(url, param string) (string, error) {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded;charset=UTF-8",
		strings.NewReader(param))
	if err != nil {
		Error(err.Error())
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("REQUEST_TIME_OUT")
		} else {
			err = errors.New("SERVER_INNER_ERROR")
		}
		return "", err
	}

	defer resp.Body.Close()
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("SERVER_INNER_ERROR")
	}

	return string(reply), nil
}
