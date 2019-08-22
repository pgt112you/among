package conn

import (
	"fmt"
	"net"
)

type DBInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname,omitempty"`
}

type MySQLDBInfo struct {
	Master     []DBInfo `json:"master"`
	Slave      []DBInfo `json:"slave"`
	MasterMode int      `json:"mastermode"`
	SlaveMode  int      `json:"slavemode"`
}

func (dbi DBInfo) CreateMySQLConn() *net.Conn {
	myaddr := fmt.Sprintf("%s:%d", dbi.Host, dbi.Port)
	conn, err := net.Dial("tcp", myaddr)
	if err != nil {
		fmt.Printf("connect to mysql %s error %s", myaddr, err)
		return nil
	}
	return &conn
}
