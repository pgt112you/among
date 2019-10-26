package main

import (
	"fmt"
	"net"

	"github.com/pgt112you/among/config"
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
	allSrv := dbobj.GetAllServer()
	if allSrv == nil {
		fmt.Println("get all server err", err)
		return
	}
	fmt.Printf("all server is %v\n", allSrv)

	for _, srv := range allSrv {
		dbKey := (*srv).GetDBPath()
		dbconf := dbobj.GetDBConf(dbKey)
		if dbconf == nil {
			fmt.Printf("get db %s conf err\n", dbKey)
			continue
		}
		go (*srv).RunServer(dbconf)
	}

	ln, err := net.Listen("tcp", ":9080")
	if err != nil {
		// handle error
	}
	for {
		_, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
	}
}
