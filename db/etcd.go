package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	//"github.com/coreos/etcd/etcdserver/etcdserverpb"
	//"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/conn"
)

type AmongDB struct {
	EC *clientv3.Client
}

func NewAmongDB(ac *config.AmongConfig) (*AmongDB, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   ac.ETCDAddr,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("etcd client error", err)
		return nil, err
	}

	adb := new(AmongDB)
	adb.EC = cli
	return adb, nil
}

func (adb *AmongDB) GetAllDBConf() {
	//resp, err := (*adb.EC).Get(context.TODO(), "/among/dbserver")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := (*adb.EC).Get(ctx, "127.0.0.1:13307")
	if err != nil {
		fmt.Println("get /among/dbserver error", err)
		return
	}
	fmt.Printf("resp is %v\n", resp)
	for _, kv := range resp.Kvs {
		fmt.Println(kv.String())
	}
}

func (adb *AmongDB) GetDBConf(key string) *conn.MySQLDBInfo {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := (*adb.EC).Get(ctx, key)
	if err != nil {
		fmt.Printf("etcd get %s error %v\n", key, err)
		return nil
	}
	fmt.Printf("resp is %v\n", resp)
	if resp.Count <= 0 {
		return nil
	}
	kv := resp.Kvs[0]
	var myinfo conn.MySQLDBInfo
	err = json.Unmarshal(kv.Value, &myinfo)
	if err != nil {
		fmt.Printf("Unmarshal mysqldbinfo err %v\n", err)
		return nil
	}
	return &myinfo
}
