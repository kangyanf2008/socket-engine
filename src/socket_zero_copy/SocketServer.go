package socket_zero_copy

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
	//监听端口
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{Port: s.Port, IP: s.Ip})
	defer listen.Close()
	if err != nil {
		panic("SoketServer Start err=[" + err.Error() + "],ip=" + fmt.Sprint(s.Ip) + ",port=" + strconv.FormatInt(int64(s.Port), 10))
		return err
	}
	s.listen = listen
	for ;!s.isStop; {
		conn, err := listen.AcceptTCP() //接收客户端连接
		if err != nil {
			continue
		}
		if err == nil {
			fmt.Printf(conn.RemoteAddr().String() + " join \n")
		}
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
	if header.Event==1 {  //需要放进连接池
		//假如为被叫终端设备，则放入长连接池
		ch := NewConnectHandler( conn)
		ch.srcIp = conn.RemoteAddr().String() //连接源IP
		ch.connectId = string(header.Data)
		//连接存放连接池
		addConnectPool(ch.connectId, ch)
	} else if header.Event==2 { //主动呼叫设备进行播号转发
		connectHandler := getConnectByConnId(string(header.Data))
		if connectHandler != nil {
			connectHandler.isProxy = true
			connectHandler.proxy(conn)
			//file1,_ := os.Open("aa.txt")
			//server(conn, file1, connectHandler.Conn)
		}
	} else { //关闭连接
		conn.Close()
	}

	return true
}
