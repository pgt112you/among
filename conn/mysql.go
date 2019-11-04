package conn

import (
	"fmt"
	"net"
)

type MySQLDBConf struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Dbname   string `json:"dbname,omitempty"`
}

type SrvMySQLDBConf struct {
	Master     []*MySQLDBConf `json:"master"`
	Slave      []*MySQLDBConf `json:"slave,omitempty"`
	MasterMode int            `json:"mastermode,omitempty"`
	SlaveMode  int            `json:"slavemode,omitempty"`
}

func NewMySQLDBConf(host string, port int) *MySQLDBConf {
	mc := new(MySQLDBConf)
	mc.Host = host
	mc.Port = port
	return mc
}

func NewSrvMySQLDBConf(mc *MySQLDBConf) *SrvMySQLDBConf {
	smc := new(SrvMySQLDBConf)
	smc.Master = append(smc.Master, mc)
	return smc
}

//func (self *SrvMySQLDBConf) GetLinkBackEnd() *MySQLDBConf {
func (self *SrvMySQLDBConf) GetLinkBackEnd() BackEndPoint {
	return self.Master[0]
}

func (self *MySQLDBConf) CreateConn() net.Conn {
	myaddr := fmt.Sprintf("%s:%d", self.Host, self.Port)
	conn, err := net.Dial("tcp", myaddr)
	if err != nil {
		fmt.Printf("connect to mysql %s error %s", myaddr, err)
		return nil
	}
	return conn
}
