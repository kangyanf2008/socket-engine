package socket

import (
	"sync"
)

var pools sync.Map

/**
把连接放到连接池中
 */
func addConnectPool(connId string, conn *ConnectHandler)  {
	pools.Store(connId, conn)
}

//删除连接
func delConnectByConnId(connId string, srcAddr string) {
	conn := getConnectByConnId(connId)
	if conn != nil && conn.srcIp == srcAddr {
		pools.Delete(connId)
	}
}

//获取连接
func getConnectByConnId(connId string) *ConnectHandler{
	connectHandler, ok:= pools.Load(connId)
	if ok {
		return connectHandler.(*ConnectHandler)
	}
	return nil
}

//断开连接
func CloseAllConnect(){
	pools.Range(func(k, connectHandler interface{})bool {
		connectHandler.(*ConnectHandler).CloseConn()
		return true
	})
}