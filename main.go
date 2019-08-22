package main

import (
	"fmt"
	"net"

	"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/conn"
	"github.com/pgt112you/among/db"
)

func main() {
	ac := config.NewAmongConfig("./among.yaml")
	if ac == nil {
		fmt.Println("ac is nil")
		return
	}
	dbobj, err := db.NewAmongDB(ac)
	if err != nil {
		fmt.Println("create db object err", err)
		return
	}
	/////////// for test //////////////
	key := "127.0.0.1:13307"
	dbconf := dbobj.GetDBConf(key)
	if dbconf == nil {
		fmt.Println("get db conf err")
		return
	}
	////////// for test end ///////////

	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		// handle error
	}
	for {
		cConn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		ct := conn.NewConnTeam(&cConn, dbconf)
		if ct == nil {
			fmt.Println("new connteam err")
			continue
		}
		go ct.Run()
	}
}
