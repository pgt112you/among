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
	return srv
}

func createMySQLServerConf() *MySQLServerConf {
	return new(MySQLServerConf)
}

func (conf *MySQLServerConf) GetType() int {
	return conf.Type
}

func (self *MySQLServerConf) Unmarshal(content []byte) error {
	err := json.Unmarshal(content, &self)
	fmt.Println("err is ", err)
	fmt.Println(self.SrvPort)
	return err
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
	self.sln.Close()
	self.mark = 1
	for _, ct := range self.cts {
		ct.Close()
	}
}

func (self *MySQLServer) Reload() {
	self.Stop()
	self.cts = self.cts[:0]
	go self.Run()
}

func CreateMySQLConf(srv *MySQLServer) *conn.SrvMySQLDBConf {
	mc := conn.NewMySQLDBConf(srv.Host, srv.Port)
	smc := conn.NewSrvMySQLDBConf(mc)
	return smc
}

func (self *MySQLServer) Run() {
	if self.mark != 1 {
		fmt.Println("server status mark is not right")
	}
	if self.Proto == "" {
		self.Proto = "tcp"
	}
	lnAddr := fmt.Sprintf("0.0.0.0:%d", self.SrvPort)
	fmt.Println(lnAddr)
	ln, err := net.Listen(self.Proto, lnAddr)
	if err != nil {
		fmt.Printf("listen %s err %v\n", lnAddr, err)
		return
		// handle error
	}
	fmt.Println("aaaaaaaaaa")
	self.sln = ln
	fmt.Println("bbbbbbbbbb")
	self.mark = 2
	for {
		if self.mark != 2 {
			break
		}
		cConn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		fmt.Println("local addr is", cConn.LocalAddr().String())
		smc := CreateMySQLConf(self)
		ct := conn.NewConnTeam(cConn, smc)
		if ct == nil {
			fmt.Println("new connteam err")
			continue
		}
		go ct.Run()
		self.cts = append(self.cts, ct)
	}

}
