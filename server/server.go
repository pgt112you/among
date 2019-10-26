package server

import (
	"fmt"
	"net"

	//"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/conn"
	//"github.com/pgt112you/among/db"

	//"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/clientv3"
)

var ServerPath string = "among/server"
//var DBServerPath string = "among/dbserver"

type ServerInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Server struct {
	ServerInfo
	EV chan *clientv3.Event
}

//func (srv ServerInfo) RunServer(dbobj *db.AmongDB) {
func (srv Server) RunServer(dbconf conn.BackEndInfo) {
	lnAddr := fmt.Sprintf(":%d", srv.Port)
	ln, err := net.Listen("tcp", lnAddr)
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

func (srv Server) WatchDBServer() {
	for {
		ev, ok := <-srv.EV
		if !ok {
			fmt.Printf("%s:%d get from srv ev err\n", srv.Host, srv.Port)
			continue
		}
		if ev.Type == "PUT" {
			if ev.Kv.Value == ev.PrevKv.Value {
				continue
			}

		} else if ev.Type == "DELETE" {
			fmt.Println("delete key ", ev.Kv.Key)
		}
	}
}

func (srv Server) GetDBPath() string {
	dbKey := fmt.Sprintf("%s/%s:%d", conn.DBPath, srv.Host, srv.Port)
	return dbKey
}
