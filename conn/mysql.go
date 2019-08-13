package conn

import (
	"fmt"
	"net"
)

func createMySQLConn() *net.Conn {
	myaddr := "127.0.0.1:13307"
	conn, err := net.Dial("tcp", myaddr)
	if err != nil {
		fmt.Printf("connect to mysql %s error %s", myaddr, err)
		return nil
	}
	return &conn
}
