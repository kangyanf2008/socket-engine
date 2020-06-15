package socket_zero_copy

import (
	"acln.ro/zerocopy"
	"net"
	"os"
)
type ConnectHandler struct {
	connectId    string
	srcIp        string
	Conn         *net.TCPConn
	quit         chan bool
	isStop       bool      //应用是否退出
	isProxy      bool      //是否代理
}

//连接处理请求类
func NewConnectHandler(conn *net.TCPConn) *ConnectHandler {
	return &ConnectHandler{
		Conn: conn,
		quit:         make(chan bool, 2),
	}
}

func (c *ConnectHandler) proxy(client net.Conn) error {
	defer func() {
		c.isProxy = false
	}()
	//client数据拷贝到conn下
	go zerocopy.Transfer(c.Conn, client)
	//conn下数据拷贝到client下
	go zerocopy.Transfer(client, c.Conn)
	return nil
}

func (c *ConnectHandler) close(client net.Conn) error {
	c.isStop = true
	c.Conn.Close()
	c.isProxy=true
	return nil
}

//关闭连接处理
func (c *ConnectHandler) CloseConn() {
	c.isStop = true
	c.quit <- true
	c.quit <- true
	c.Conn.Close()
	c.isProxy=true
	delConnectByConnId(c.connectId, c.srcIp)
	close(c.quit)
}



func proxy(upstream, downstream net.Conn) error {
	//up数据拷贝到down下
	go zerocopy.Transfer(downstream, upstream)
	//down下数据拷贝到down下
	go zerocopy.Transfer(upstream, downstream)
	return nil

}

func server(camera net.Conn, recording *os.File, client net.Conn)  error {
	campipe, err := zerocopy.NewPipe()
	if err != nil {
		return err
	}
	defer campipe.Close()

	// Create a pipe to the recording.
	recpipe, err := zerocopy.NewPipe()
	if err != nil {
		return err
	}
	defer recpipe.Close()

	// Arrange for data on campipe to be duplicated to recpipe.
	campipe.Tee(recpipe)

	// Run the world.
	go campipe.ReadFrom(camera)
	go recpipe.WriteTo(recording)
	go campipe.WriteTo(client)
	return nil
}