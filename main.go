package main

import (
	"fmt"
	"net"

	"github.com/pgt112you/among/conn"
)

func main() {
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
		ct := conn.NewConnTeam(&cConn)
		if ct == nil {
			fmt.Println("new connteam err")
			continue
		}
		go ct.Run()
	}
}
