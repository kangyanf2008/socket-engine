package socket

import (
	"fmt"
	"net"
	"time"
)


type ConnectHandler struct {
	connectId    string
	srcIp        string
	Conn         *net.TCPConn
	quit         chan bool
	RequestQueue chan *Header
	writeQueue   chan []byte
	isStop       bool
	dealMsgNum   int64
}

//连接处理请求类
func NewConnectHandler(writeQueueLen int, requestQueueLen int, conn *net.TCPConn) *ConnectHandler {
	if writeQueueLen == 0 {
		writeQueueLen = 100
	}
	if requestQueueLen == 0 {
		requestQueueLen = 100
	}
	return &ConnectHandler{Conn: conn,
		writeQueue:   make(chan []byte, writeQueueLen),
		RequestQueue: make(chan *Header, requestQueueLen),
		quit:         make(chan bool, 3),
	}
}

//处理连接请求
func (c *ConnectHandler) Run() {
	go c.dealRequest()
	go c.read()
	go c.write()
}

//向客户端发消息
func (c *ConnectHandler) PushMessage(data []byte) {
	if !c.isStop {
		c.writeQueue <- data
	}
}

//关闭连接处理
func (c *ConnectHandler) CloseConn() {
	c.isStop = true
	c.quit <- true
	c.quit <- true
	c.quit <- true
	delConnectByConnId(c.connectId, c.srcIp)
	close(c.quit)
}


//读处理
func (c *ConnectHandler) read() {
	headerByte := make([]byte, HeaderSize)
	defer func() {
		c.Conn.CloseRead()
		c.CloseConn()
	}()
	for ; !c.isStop; {
		readSize, err := c.Conn.Read(headerByte)
		if err != nil {
			return
		}
		if readSize != HeaderSize {
			continue
		}
		//解析协议头和消息体
		header := unPackHeader(headerByte, c.Conn)
		if header != nil {
			//请求数据放队列进行异步处理
			c.RequestQueue <- header
		}
	}
}


//写处理
func (c *ConnectHandler) write() {
	isStop := false
	defer func() {
		//close(c.writeQueue)
		c.Conn.CloseWrite()
	}()

	for {
		select {
		case <-c.quit:
			isStop = true
		case w := <-c.writeQueue:
			_, err := c.Conn.Write(w)
			if err != nil {
				fmt.Printf("writeQueue write goroutine exit connectId=[%s], srcIp=[%s] cause error[%s] \n", c.connectId, c.srcIp, err)
				return
			}
		default: //如果超过100毫秒没有数据，并且处理连接为关闭状态，则关闭写队列
			time.Sleep(100 * time.Millisecond)
			if isStop {
fmt.Printf("write goroutine exit connectId=[%s], srcIp=[%s] cause isStop disconnect \n", c.connectId, c.srcIp)
				return
			}
		}
	}

}

//处理请求
func (c *ConnectHandler) dealRequest() {
	isStop := false
	defer func() {
		c.Conn.CloseRead()
		//close(c.RequestQueue)
	}()
	for {
		select {
		case <-c.quit:
			isStop = true
		case req := <-c.RequestQueue:
			c.dealMsgNum ++
			if c.dealMsgNum % 10000 == 0 {
				if req != nil {
					fmt.Println(time.Now().Unix(),",", c.dealMsgNum)
				}
			}
			//写回客户端
			c.PushMessage(packHeader(req))
			//fmt.Printf("%s ==\n", req.Data)
		default:
			time.Sleep(100 * time.Millisecond) //超过100毫秒没有数据，并且连接为关闭状态，则关闭请求队列
			if isStop {
				return
			}
		}
	}
}
