package db

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/coreos/etcd/clientv3"
	//"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/conn"
	"github.com/pgt112you/among/server"
)

type AmongDB struct {
	EC *clientv3.Client
}

const (
	PUTOP = mvccpb.PUT
	DELOP = mvccpb.DELETE
)

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

func (adb *AmongDB) GetAllServerConf() []server.CommonServerConf {
	//resp, err := (*adb.EC).Get(context.TODO(), "/among/dbserver")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	opts := clientv3.WithPrefix()
	resp, err := (*adb.EC).Get(ctx, server.ServerPath, opts)
	if err != nil {
		fmt.Println("get /among/dbserver error", err)
		return nil
	}
	//srvs := make([]*server.Server, len(resp.Kvs))
	srvConfs := make([]server.CommonServerConf, len(resp.Kvs))
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
		case server.MySQLSTR:
			msc := new(server.MySQLServerConf)
			fmt.Println(string(kv.Value))
			err := msc.Unmarshal(kv.Value)
			fmt.Println("--", msc.SrvPort)
			if err != nil {
				fmt.Println("unmarshal mysql serverinfo error")
				continue
			}
			msc.Type = server.MySQL
			fmt.Println("==", msc.SrvPort)
			//srv := new(server.MySQLServer)
			//srv.EV = make(chan *clientv3.Event)
			srvConfs[i] = msc

		default:
			fmt.Println("server type is error")
			continue
		}
	}
	return srvConfs
}

func (adb *AmongDB) GetServerConf(srvType, addr string) server.CommonServerConf {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	key := fmt.Sprintf("%s/%s/%s", server.ServerPath, srvType, addr)
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
	if srvType == server.MySQLSTR {
		msc := new(server.MySQLServerConf)
		err := msc.Unmarshal(kv.Value)
		if err != nil {
			fmt.Println("unmarshal mysql serverconf error")
			return nil
		}
		msc.Type = server.MySQL
		return msc
	} else {
		fmt.Println("wrong server type")
		return nil
	}
}

func (adb *AmongDB) GetDBConf(key string) *conn.MySQLDBConf {
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
	var myinfo conn.MySQLDBConf
	err = json.Unmarshal(kv.Value, &myinfo)
	if err != nil {
		fmt.Printf("Unmarshal mysqldbinfo err %v\n", err)
		return nil
	}
	return &myinfo
}

func (adb *AmongDB) WatchAllServer(evchan chan *clientv3.Event) {
	rch := (*adb.EC).Watch(context.Background(), server.ServerPath, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			evchan <- ev
			fmt.Printf("now event %s %q  %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			//fmt.Printf("prev event %s %s\n", string(ev.PrevKv.Key), string(ev.PrevKv.Value))
		}
	}
}
