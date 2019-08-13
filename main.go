package main

import (
	"fmt"
	"net"
	"time"

	"github.com/pgt112you/among/conn"
	"github.com/pgt112you/among/tool"
)

func createMySQLConn() *net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:13307")
	if err != nil {
		fmt.Println(err)
		fmt.Println("connect to mysql error", err)
		return nil
	}
	return &conn
}

func handleConnection(cliconn *net.Conn) {
	srvconn := createMySQLConn()
	if srvconn == nil {
		fmt.Println("connect to mysql error")
		return
	}
	var ct conn.ConnTeam
	ct.CliConn = cliconn
	ct.SrvConn = srvconn
	ct.CliDataChan = make(chan []byte)
	ct.SrvDataChan = make(chan []byte)
	go dealSrv(&ct)
	time.Sleep(100 * time.Millisecond)
	go dealCli(&ct)
}

func srvSend(ct *conn.ConnTeam) {
	for {
		data := <-ct.SrvDataChan
		packLen := len(data)
		fmt.Println("in srvsend len is", packLen)
		n := 0
		for n < packLen {
			tempN, err := (*ct.SrvConn).Write(data[n:])
			if err != nil {
				return
			}
			n += tempN
			fmt.Printf("in srvsend n is %d\n", n)
		}
	}
}

func dealSrv(ct *conn.ConnTeam) {
	go srvSend(ct)
	for {
		fmt.Println("ddddddddddddddddd")
		data := tool.GetData(*ct.SrvConn, "server")
		fmt.Println("DDDDDDDDDDDDDDDDD", len(data))
		if data == nil {
			fmt.Println("in dealsrv data is nil")
			return
		}
		ct.CliDataChan <- data
	}
}

func cliSend(ct *conn.ConnTeam) {
	for {
		data := <-ct.CliDataChan
		packLen := len(data)
		n := 0
		fmt.Println("in clisend len is", packLen)
		for n < packLen {
			tempN, err := (*ct.CliConn).Write(data[n:])
			if err != nil {
				return
			}
			n += tempN
			fmt.Printf("in clisend n is %d\n", n)
		}
	}
}

func dealCli(ct *conn.ConnTeam) {
	go cliSend(ct)
	for {
		fmt.Println("ccccccccccccccccc")
		data := tool.GetData(*ct.CliConn, "client")
		fmt.Println("CCCCCCCCCCCCCCCCC", len(data))
		if data == nil {
			fmt.Println("in dealcli data is nil")
			return
		}
		ct.SrvDataChan <- data
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(&conn)
	}
}
