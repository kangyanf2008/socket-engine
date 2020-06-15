package socket

import (
	"fmt"
	"net"
	"strconv"
)

type SoketServer struct {
	Port   int
	Ip     net.IP
	isStop bool
	listen *net.TCPListener
}

func NewSoketServer(port int) *SoketServer {
	if port == 0 {
		port = 8080
	}
	return &SoketServer{Port:port}
}

func (s *SoketServer) Start() error {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{Port: s.Port, IP: s.Ip})
	defer listen.Close()
	if err != nil {
		panic("SoketServer Start err=[" + err.Error() + "],ip=" + fmt.Sprint(s.Ip) + ",port=" + strconv.FormatInt(int64(s.Port), 10))
		return err
	}
	s.listen = listen
	for ;!s.isStop; {
		conn, err := listen.AcceptTCP()
		if err != nil {
			continue
		}
		fmt.Printf(conn.RemoteAddr().String() + " join \n")
		//连接认证通过后，进行逻辑处理
		if !s.check(conn) {
			conn.Close()
		}
	}
	return nil
}

func (s *SoketServer) Stop() {
	s.isStop = true
	s.listen.Close()
	//关闭所有client连接
	CloseAllConnect()
}

//连接认证
func (s *SoketServer) check(conn *net.TCPConn) bool {
	buffer := make([]byte, HeaderSize)
	readSize, err := conn.Read(buffer)
	if err != nil || readSize != HeaderSize {
		return false
	}
	header := unPackHeader(buffer, conn)
	if header == nil {
		return false
	}
	//启动连接处理协程
	ch := NewConnectHandler(50,50, conn,)
	ch.srcIp = conn.RemoteAddr().String() //连接源IP
	ch.connectId = string(header.Data)
	ch.Run()
	//连接存放连接池
	addConnectPool(ch.connectId, ch)

	return true
}
