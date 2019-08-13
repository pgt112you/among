package conn

import (
	"net"
)

type ConnTeam struct {
	CliConn     *net.Conn
	SrvConn     *net.Conn
	CliDataChan chan []byte
	SrvDataChan chan []byte
}
