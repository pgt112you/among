package server

import (
	"fmt"
)

const (
	MySQL = 1
)

const (
	MySQLSTR = "MySQL"
)

const ServerDBPath = "among/server"

type CommonServerConf interface {
	Unmarshal([]byte) error
	GetType() int
}

type CommonServer interface {
	Run()
	Stop()
	Reload()
	CheckOk() bool
	GetAddr() string
	SetSrvPort(int)
}

func CreateSrv(srvConf CommonServerConf) CommonServer {
	switch srvConf.GetType() {
	case MySQL:
		mc, ok := srvConf.(*MySQLServerConf)
		if !ok {
			return nil
		}
		return createMySQLServer(mc)
	default:
		mc, ok := srvConf.(*MySQLServerConf)
		if !ok {
			return nil
		}
		return createMySQLServer(mc)
	}
}

func CreateSrvConf(ty string, confContent []byte) CommonServerConf {
	switch ty {
	case MySQLSTR:
		srvc := new(MySQLServerConf)
		if confContent == nil {
			return srvc
		}
		err := srvc.Unmarshal(confContent)
		if err != nil {
			fmt.Println("unmarsh mysql conf error", err)
			return nil
		}
		return srvc
	default:
		srvc := new(MySQLServerConf)
		if confContent == nil {
			return srvc
		}
		err := srvc.Unmarshal(confContent)
		if err != nil {
			fmt.Println("unmarsh mysql conf error", err)
			return nil
		}
		return srvc
	}
}
