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
	"github.com/pgt112you/among/server"
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

func (adb *AmongDB) GetAllServer() []*server.Server {
	//resp, err := (*adb.EC).Get(context.TODO(), "/among/dbserver")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	opts := clientv3.WithPrefix()
	resp, err := (*adb.EC).Get(ctx, server.ServerPath, opts)
	if err != nil {
		fmt.Println("get /among/dbserver error", err)
		return nil
	}
	//srvs := make([]*server.Server, len(resp.Kvs))
	fmt.Printf("resp.kvs len is %d\n", len(resp.Kvs))
	srvs := make([]*server.Server, len(resp.Kvs))
	fmt.Printf("all server is %v\n", srvs)
	for i, kv := range resp.Kvs {
		fmt.Println(kv.Key, kv.Value)
		var si server.ServerInfo
		err = json.Unmarshal(kv.Value, &si)
		if err != nil {
			fmt.Printf("Unmarshal server info err %v\n", err)
			return nil
		}
		srv := new(server.Server)
		srv.ServerInfo = si
		srv.EV = make(chan *clientv3.Event)
		fmt.Printf("srv is %v\n", srv)
		//srvs = append(srvs, srv)
		srvs[i] = srv
	}
	fmt.Printf("all server is %v\n", srvs)
	return srvs
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

//func (adb *AmongDB) WatchServer(srvKey string) {
func (adb *AmongDB) WatchServer(srv *server.Server) {
	srvKey := fmt.Sprintf("%s/%s:%d", server.DBServerPath, srv.Host, srv.Port)
	rch := (*adb.EC).Watch(context.Background(), srvKey, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			srv.EV <- ev
			fmt.Printf("%s %q  %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			fmt.Printf("%s %s\n", string(ev.PrevKv.Key), string(ev.PrevKv.Value))
		}
	}
}
