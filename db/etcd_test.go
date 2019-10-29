package db

import (
	"fmt"
	"testing"

	"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/db"
)

func TestMain(t *testing.T) {
	ac := config.NewAmongConfig("../among.yaml")
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
}
