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

type MySQLBackEndConf struct {
	Master     []*MySQLDBConf `json:"master"`
	Slave      []*MySQLDBConf `json:"slave,omitempty"`
	MasterMode int            `json:"mastermode,omitempty"`
	SlaveMode  int            `json:"slavemode,omitempty"`
}

type MySQLBackEnd struct {
	MySQLBackEndConf
	MasterConn []net.Conn
}

func NewMySQLDBConf(host string, port int) *MySQLDBConf {
	mc := new(MySQLDBConf)
	mc.Host = host
	mc.Port = port
	return mc
}

func NewMySQLBackEndConf(mc *MySQLDBConf) *MySQLBackEndConf {
	mbec := new(MySQLBackEndConf)
	mbec.Master = append(mbec.Master, mc)
	return mbec
}

func createMySQLBackEnd(mbec *MySQLBackEndConf) *MySQLBackEnd {
	mbe := new(MySQLBackEnd)
	mbe.MySQLBackEndConf = *mbec
	//mbe.MasterConn = make([]net.Conn, 1)
	return mbe
}

func (self *MySQLBackEndConf) GetType() ConfType {
	return MySQL
}

func (self *MySQLBackEnd) CreateConn() error {
	mc := self.Master[0]
	addr := fmt.Sprintf("%s:%d", mc.Host, mc.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("create mysql backend conn error", err)
		return err
	}
	self.MasterConn = append(self.MasterConn, conn)
	return nil
}

func (self *MySQLBackEnd) SendData(data []byte) (int, error) {
	return sendData(self.MasterConn[0], data)
}

func (self *MySQLBackEnd) RecvData() ([]byte, int, error) {
	return recvData(self.MasterConn[0])
}

func (self *MySQLBackEnd) Close() {
	for _, conn := range self.MasterConn {
		conn.Close()
	}
}
