package socket

import (
	"bytes"
	"encoding/binary"
	"net"
)

const HeaderSize int = 5

type Header struct {
	Event   byte   //事件
	DataLen uint32  //数据长度
	Data    []byte //内容
}

//解包
func unPackHeader(headerByte []byte, conn *net.TCPConn) *Header{
	headerBuf := bytes.NewReader(headerByte)
	var event byte
	var dataLen uint32
	var data []byte
	binary.Read(headerBuf, binary.BigEndian, &event)
	binary.Read(headerBuf, binary.BigEndian, &dataLen)
	if dataLen > 0 {
		data = make([]byte, dataLen)
		readSize, err := conn.Read(data)
		if uint32(readSize) != dataLen || err != nil {
			return nil
		}
	}
	return &Header{Event: event, DataLen: dataLen, Data: data}
}

//封包
func packHeader(h *Header) []byte {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(h.Event)
	binary.Write(buffer, binary.BigEndian, uint32(h.DataLen))
	buffer.Write(h.Data)
	return buffer.Bytes()
}