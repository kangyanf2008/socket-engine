package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"socket"
	"time"
)

func main() {
	go func() {
		target, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", "192.168.34.114", 10000))
		//target, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", "192.168.26.128", 10000))
		conn, err := net.DialTCP("tcp4", nil, target)
		checkError2(err)

		//event1进行连接注册
		data := "101"
		databytes := []byte(data)
		dataLen := len(databytes)
		buffer := bytes.NewBuffer(nil)
		buffer.WriteByte(1)
		binary.Write(buffer, binary.BigEndian, uint32(dataLen))
		buffer.Write(databytes)
		_, err = conn.Write(buffer.Bytes())
fmt.Println("login success", err)
		isFirst := true
		//读数据
		headerByte := make([]byte, socket.HeaderSize)

		donwMsgNum := 0  //下行消息数量
		for {
			readSize, err := conn.Read(headerByte)
			if err != nil {
				break
			}
			if readSize != socket.HeaderSize {
				continue
			}
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
					continue
				}
//fmt.Println(string(data))    //打印转发内容
donwMsgNum++
if donwMsgNum % 10000 == 0 {
		fmt.Println(time.Now().Unix(),",", donwMsgNum)
}
			}
			//收到第一条消息后，进行异步写消息
			if isFirst {
				isFirst = false
				//异常写数据
				go func() {
					data := "333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333"
					databytes := []byte(data)
					for ; ; {
						dataLen := len(databytes)
						buffer := bytes.NewBuffer(nil)
						buffer.WriteByte(1)
						binary.Write(buffer, binary.BigEndian, uint32(dataLen))
						buffer.Write(databytes)
						_, err = conn.Write(buffer.Bytes())
//fmt.Println(data)
						checkError2(err)
					}
				}()
			}

		} // end for

	}()
	//os.Exit(0)
	select {}
}
func checkError2(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
