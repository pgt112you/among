package conn

import (
	"fmt"
)

type ConfType int

const (
	MySQL ConfType = 1
)


type BackEnd interface {
	CreateConn() error
	SendData([]byte) (int, error)
	RecvData() ([]byte, int, error)
	Close()
}

type BackEndConf interface {
    GetType() ConfType
}

func createBackEnd(bec BackEndConf) (BackEnd, error){
	switch bec.GetType() {
	case MySQL:
		mbec, ok := bec.(*MySQLBackEndConf)
		if !ok {
			return nil, fmt.Errorf("config type invalid")
		}
		mbe := createMySQLBackEnd(mbec)
		if mbe != nil {
			return mbe, nil
		} else {
			return nil, fmt.Errorf("create mysql backend error")
		}
	default:
		mbec, ok := bec.(*MySQLBackEndConf)
		if !ok {
			return nil, fmt.Errorf("config type invalid")
		}
		mbe := createMySQLBackEnd(mbec)
		if mbe != nil {
			return mbe, nil
		} else {
			return nil, fmt.Errorf("create mysql backend error")
		}
	}

}
