package conn

import (
	"fmt"
	"io"
	"net"
	"time"
)

type ConnTeam struct {
	CliConn     net.Conn
	SrvConn     net.Conn
	CliDataChan chan []byte
	SrvDataChan chan []byte
	CliStatus   int
	SrvStatus   int
}

func NewConnTeam(cConn net.Conn, beinfo BackEndInfo) *ConnTeam {
	dbinfo := beinfo.GetLinkBackEnd()
	sConn := dbinfo.CreateConn()
	if sConn == nil {
		fmt.Println("connect to mysql error")
		return nil
	}
	var ct ConnTeam
	ct.CliConn = cConn
	ct.SrvConn = sConn
	ct.CliDataChan = make(chan []byte, 10)
	ct.SrvDataChan = make(chan []byte, 10)
	ct.CliStatus = 0
	ct.SrvStatus = 0
	return &ct
}

func (ct *ConnTeam) Close() {
	ct.CloseConn(2)
}

func (ct *ConnTeam) CloseConn(who int) { // who 1 client, who 2 server
	if who == 1 {
		//(*ct.CliConn).Close()
		ct.CliConn.Close()
		ct.CliStatus = -1
		time.Sleep(100 * time.Millisecond)
		//(*ct.SrvConn).Close()
		ct.SrvConn.Close()
		ct.SrvStatus = -1
	} else {
		fmt.Println("server close")
		//(*ct.SrvConn).Close()
		ct.SrvConn.Close()
		ct.SrvStatus = -1
		time.Sleep(100 * time.Millisecond)
		//(*ct.CliConn).Close()
		ct.CliConn.Close()
		ct.CliStatus = -1
	}
}

func (ct *ConnTeam) Run() {
	go ct.dealSrv()
	time.Sleep(100 * time.Millisecond)
	go ct.dealCli()
	for ct.CliStatus != -1 && ct.SrvStatus != -1 {
		time.Sleep(1 * time.Second)
	}

}

func (ct *ConnTeam) srvSend() {
	for ct.SrvStatus != -1 {
		data := <-ct.SrvDataChan
		packLen := len(data)
		n, err := sendData(ct.SrvConn, data)
		if err != nil {
			fmt.Printf("srvsend send %d less than %d\n", n, packLen)
		}
	}
}

func (ct *ConnTeam) dealSrv() {
	go ct.srvSend()
	for ct.SrvStatus != -1 {
		data, _, err := recvData(ct.SrvConn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("server close, close connection")
			} else {
				fmt.Println("recv server data err", err)
			}
			ct.CloseConn(2)
			return
		}
		ct.CliDataChan <- data
	}
}

func (ct *ConnTeam) cliSend() {
	for ct.CliStatus != 1 {
		data := <-ct.CliDataChan
		packLen := len(data)
		n, err := sendData(ct.CliConn, data)
		if err != nil {
			fmt.Printf("clisend send %d less than %d\n", n, packLen)
		}
	}
}

func (ct *ConnTeam) dealCli() {
	go ct.cliSend()
	for ct.CliStatus != -1 {
		data, _, err := recvData(ct.CliConn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client close, close connection")
			} else {
				fmt.Println("recv client data err", err)
			}
			ct.CloseConn(1)
			return
		}
		ct.SrvDataChan <- data
	}
}
