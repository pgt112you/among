package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/coreos/etcd/clientv3"
	"github.com/pgt112you/among/conn"
)

type MySQLServerConf struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	SrvPort int    `json:"srvport"`
	Proto   string `json:"protocol,omitempty"`
	Type    int    `json:"type,omitempty"`
}

type MySQLServer struct {
	MySQLServerConf
	EV   chan *clientv3.Event
	mark int
	sln  net.Listener
	cts  []*conn.ConnTeam
}

func createMySQLServer(conf *MySQLServerConf) *MySQLServer {
	srv := new(MySQLServer)
	srv.MySQLServerConf = *conf
	srv.EV = make(chan *clientv3.Event)
	srv.mark = 1
	if srv.Proto == "" {
		srv.Proto = "tcp"
	}
	lnAddr := fmt.Sprintf("0.0.0.0:%d", srv.SrvPort)
	ln, err := net.Listen(srv.Proto, lnAddr)
	if err != nil {
		fmt.Printf("listen %s err %v\n", lnAddr, err)
		return nil
		// handle error
	}
	srv.sln = ln
	return srv
}

func createMySQLServerConf() *MySQLServerConf {
	return new(MySQLServerConf)
}

func CreateMySQLConf(srv *MySQLServer) *conn.MySQLBackEndConf {
	mc := conn.NewMySQLDBConf(srv.Host, srv.Port)
	smc := conn.NewMySQLBackEndConf(mc)
	return smc
}

func (conf *MySQLServerConf) GetType() int {
	return conf.Type
}

func (self *MySQLServerConf) Unmarshal(content []byte) error {
	//err := json.Unmarshal(content, &self)
	err := json.Unmarshal(content, self)
	return err
}

//func (self *MySQLServerConf) GetBackConf() conn.BackEndConf {
func (self *MySQLServerConf) GetBackConf() *conn.MySQLBackEndConf {
	//myc := new(conn.MySQLDBConf)
	mbc := new(conn.MySQLBackEndConf)
	mc := new(conn.MySQLDBConf)
	mc.Host = self.Host
	mc.Port = self.Port
	mbc.Master = make([]*conn.MySQLDBConf, 1)
	mbc.Master[0] = mc
	fmt.Println(mbc)
	return mbc
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *MySQLServer) CompareConf(dstC CommonServerConf) bool {
	dst, ok := dstC.(*MySQLServerConf)
	if !ok {
		return false
	}
	if self.Host != dst.Host {
		return false
	}
	if self.Port != dst.Port {
		return false
	}
	if self.SrvPort != dst.SrvPort {
		return false
	}
	return true
}

func (self *MySQLServer) SetSrvPort(port int) {
	self.SrvPort = port
}

func (self *MySQLServer) GetAddr() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *MySQLServer) CheckOk() bool {
	if self.mark == 2 {
		return true
	} else {
		return false
	}
}

func (self *MySQLServer) Stop() {
	self.mark = 1
	self.sln.Close()
	for _, ct := range self.cts {
		ct.Close()
	}
}

func (self *MySQLServer) Reload() {
	self.Stop()
	self.cts = self.cts[:0]
	go self.Run()
}

func (self *MySQLServer) Run() {
	if self.mark != 1 {
		fmt.Println("server status mark is not right")
		return
	}
	self.mark = 2
	for {
		if self.mark != 2 {
			break
		}
		cConn, err := self.sln.Accept()
		if err != nil {
			closeErr := fmt.Sprintf("accept tcp [::]:%d: use of closed network connection", self.SrvPort)
			if err.Error() == closeErr {
				fmt.Println("ln conn is closed >>>>>>>>>>>>>>>>>>>>>>>>>>>")
				self.sln.Close()
				self.mark = 1
				break
			}
			// handle error
			continue
		}
		fmt.Println("local addr is", cConn.LocalAddr().String())
		mbc := self.GetBackConf()
		ct := conn.NewConnTeam(cConn, mbc)
		if ct == nil {
			fmt.Println("new connteam err")
			continue
		}
		go ct.Run()
		self.cts = append(self.cts, ct)
	}

}
