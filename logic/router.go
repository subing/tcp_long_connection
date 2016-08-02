package logic

import (
	"errors"
)

func Router(cmd string, data map[string]interface{}) (string, error) {
	switch cmd {
	case "login":
		return Login(data)
	default:
		return "", errors.New("bad cmd")
	}
}
