package tcpserver

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"tcp_long_connection/logger"
)

type Package struct {
	Length  int32
	Version int32
	Flag    int32
	Serial  int32
}

func (p *Package) SetPackage(version, flag, serial int) {
	p.Version = int32(version)
	p.Flag = int32(flag)
	p.Serial = int32(serial)
}
func (p *Package) GetHeader() Package {
	return *p
}
func (p *Package) Encode(message string) ([]byte, error) {
	p.Length = int32(len(message))
	// 读取消息的长度
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息长度
	err := binary.Write(pkg, binary.LittleEndian, p.Length)
	if err != nil {
		return nil, err
	}
	//写消息版本
	err = binary.Write(pkg, binary.LittleEndian, p.Version)
	if err != nil {
		return nil, err
	}
	//写消息类型
	err = binary.Write(pkg, binary.LittleEndian, p.Flag)
	if err != nil {
		return nil, err
	}
	//写消息序号
	err = binary.Write(pkg, binary.LittleEndian, p.Serial)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}
func Decode(reader *bufio.Reader) (Package, string, error) {
	var p Package
	// 读取消息的长度
	lengthByte := make([]byte, 16)
	n, err := reader.Read(lengthByte)
	if err != nil {
		return p, "", err
	}
	if n < 16 {
		logger.Error("header read failed header not 16 byte")
		return p, "", err
	}
	lengthBuff := bytes.NewBuffer(lengthByte[:4])
	logger.Info("lengthByte:=", lengthByte, " lengthBuff:= ", lengthBuff)
	err = binary.Read(lengthBuff, binary.LittleEndian, &p.Length)
	if err != nil {
		return p, "", err
	}
	//读取消息版本
	versionBuff := bytes.NewBuffer(lengthByte[4:8])
	err = binary.Read(versionBuff, binary.LittleEndian, &p.Version)
	if err != nil {
		return p, "", err
	}
	//读取消息类型
	flagBuff := bytes.NewBuffer(lengthByte[8:12])
	err = binary.Read(flagBuff, binary.LittleEndian, &p.Flag)
	if err != nil {
		return p, "", err
	}
	//读取消息序号
	serialBuff := bytes.NewBuffer(lengthByte[12:16])
	err = binary.Read(serialBuff, binary.LittleEndian, &p.Serial)
	if err != nil {
		return p, "", err
	}
	if int32(reader.Buffered()) < p.Length {
		return p, "", err
	}
	// 读取消息真正的内容
	pack := make([]byte, int(p.Length))
	_, err = reader.Read(pack)
	if err != nil {
		return p, "", err
	}
	logger.Info("msg := ", string(pack))
	return p, string(pack), nil
}
