package conn

import (
	"net"
)

type BackEndPoint interface {
	CreateConn() net.Conn
}

type BackEndInfo interface {
	GetLinkBackEnd() BackEndPoint
}
