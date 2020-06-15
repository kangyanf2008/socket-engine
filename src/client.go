package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"socket"
)

func main() {
	for i:=0; i<200; i++ {
		go func() {
			target, _ :=net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d","192.168.34.114", 10000))
			conn, err :=net.DialTCP("tcp4",nil,target)
			checkError(err)
			//异常写数据
			go func() {
				data := "dataaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaacccNNN"
				databytes := []byte(data)
				for ; ;  {
					dataLen :=len(databytes)
					buffer := bytes.NewBuffer(nil)
					buffer.WriteByte(1)
					binary.Write(buffer, binary.BigEndian, uint32(dataLen))
					buffer.Write(databytes)
					_, err = conn.Write(buffer.Bytes())
					checkError(err)
				}
			}()
			//读数据
			headerByte := make([]byte, socket.HeaderSize)
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
//fmt.Println(string(data))
				}
			} // end for

		}()
	}
	//os.Exit(0)
	select {
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}