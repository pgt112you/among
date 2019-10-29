package db

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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

func (adb *AmongDB) GetAllServer() []*server.CommonServer {
	//resp, err := (*adb.EC).Get(context.TODO(), "/among/dbserver")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	opts := clientv3.WithPrefix()
	resp, err := (*adb.EC).Get(ctx, server.ServerPath, opts)
	if err != nil {
		fmt.Println("get /among/dbserver error", err)
		return nil
	}
	//srvs := make([]*server.Server, len(resp.Kvs))
	srvs := make([]*server.CommonServer, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		fmt.Println(string(kv.Key), string(kv.Value))
		keyStr := string(kv.Key)
		re := regexp.MustCompile(`among/server/(\w+)/(.+)`)
		reSrvInfo := re.FindAllStringSubmatch(keyStr, -1)
		if len(reSrvInfo[0]) < 3 {
			continue
		}
		sType := reSrvInfo[0][1]
		sAddr := reSrvInfo[0][2]

		fmt.Printf("stype is %s, saddr is %s\n", sType, sAddr)
		switch sType {
		case "MySQL":
			var msi server.MySQLServerInfo
			err := msi.Unmarshal(kv.Value)
			if err != nil {
				fmt.Println("unmarshal mysql serverinfo error")
			}
			srv := new(server.MySQLServer)
			srv.EV = make(chan *clientv3.Event)
			srvs[i] = srv

		default:
			fmt.Println("server type is error")
			continue
		}
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
	//srvKey := fmt.Sprintf("%s/%s:%d", server.DBServerPath, srv.Host, srv.Port)
	srvKey := fmt.Sprintf("%s/%s:%d", server.ServerPath, srv.Host, srv.Port)
	rch := (*adb.EC).Watch(context.Background(), srvKey, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			srv.EV <- ev
			fmt.Printf("now event %s %q  %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			fmt.Printf("prev event %s %s\n", string(ev.PrevKv.Key), string(ev.PrevKv.Value))
		}
	}
}
