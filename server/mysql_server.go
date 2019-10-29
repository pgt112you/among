package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/coreos/etcd/clientv3"

	"github.com/pgt112you/among/conn"
)

type mAddrInfo struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proto string `json:"protocol,omitempty"`
}

type MySQLServerInfo struct {
	mAddrInfo
	//Host string `json:"host"`
	//Port int    `json:"port"`
}

type MySQLServer struct {
	MySQLServerInfo
	EV chan *clientv3.Event
}

func (self *MySQLServerInfo) Unmarshal(content []byte) error {
	err := json.Unmarshal(content, self)
	return err
}

func (self *MySQLServerInfo) Run() {
	if self.Proto == "" {
		self.Proto = "tcp"
	}
	lnAddr := fmt.Sprintf("%s:%d", self.Host, self.Port)
	ln, err := net.Listen(self.Proto, lnAddr)
	if err != nil {
		fmt.Printf("listen %s err %v\n", lnAddr, err)
		return
		// handle error
	}
	for {
		cConn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		fmt.Println("local addr is", cConn.LocalAddr().String())
		ct := conn.NewConnTeam(&cConn, dbconf)
		if ct == nil {
			fmt.Println("new connteam err")
			continue
		}
		go ct.Run()
	}

}
