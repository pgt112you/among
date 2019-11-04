package main

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/coreos/etcd/clientv3"

	"github.com/pgt112you/among/config"
	"github.com/pgt112you/among/db"
	"github.com/pgt112you/among/server"
)

func main() {
	ac := config.NewAmongConfig("./among.yaml")
	if ac == nil {
		fmt.Println("ac is nil")
		return
	}
	dbobj, err := db.NewAmongDB(ac)
	if err != nil {
		fmt.Println("create db object err", err)
		return
	}
	allSrvConf := dbobj.GetAllServerConf()
	if allSrvConf == nil {
		fmt.Println("get all server err", err)
		return
	}
	fmt.Printf("all server is %v\n", allSrvConf)

	srvs := make(map[string]server.CommonServer)

	for _, srvc := range allSrvConf {
		// 防止重复的key加入到map中
		srv := server.CreateSrv(srvc)
		_, ok := srvs[srv.GetAddr()]
		if ok {
			fmt.Printf("%s server is  already existed", srv.GetAddr())
			continue
		}
		srvs[srv.GetAddr()] = srv
		go srv.Run()
		time.Sleep(100 * time.Millisecond)
		if !srv.CheckOk() {
			delete(srvs, srv.GetAddr())
		}
	}

	evchan := make(chan *clientv3.Event)
	go dbobj.WatchAllServer(evchan)
	ln, err := net.Listen("tcp", ":9080")
	if err != nil {
		// handle error
	}
	go dealSrvConfChange(srvs, evchan)
	go func() {
		http.ListenAndServe("0.0.0.0:6666", nil) // 这里还需要import "net/http"
	}()
	for {
		_, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
	}
}

func dealSrvConfChange(m map[string]server.CommonServer, ch chan *clientv3.Event) {
	for {
		select {
		case ev := <-ch:
			fmt.Printf("now event %s %q  %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.PrevKv != nil {
				fmt.Printf("prev event %s %s\n", string(ev.PrevKv.Key), string(ev.PrevKv.Value))
			}
			fmt.Printf("%v %s\n", ev.Type, string(ev.Type))

			re := regexp.MustCompile(`among/server/(\w+)/(.+)`)
			fmt.Println(string(ev.Kv.Key))
			res := re.FindAllStringSubmatch(string(ev.Kv.Key), -1)
			fmt.Printf("==== %v\n", res)
			if len(res[0]) < 3 {
				continue
			}
			fmt.Printf("%v\n", m)
			srv := m[res[0][2]]
			fmt.Printf("%v\n", srv)

			if ev.Type == db.PUTOP {
				if ev.PrevKv != nil { // modify existed one
					if res[0][1] == server.MySQLSTR {
						mc := new(server.MySQLServerConf)
						mc.Unmarshal(ev.Kv.Value)
						srv.SetSrvPort(mc.SrvPort)
					}
					go srv.Reload()
				} else { // add new one
					srvType := res[0][1]
					srvc := server.CreateSrvConf(srvType, ev.Kv.Value)
					if srvc == nil {
						fmt.Println("add new server conf error")
						continue
					}
					srv := server.CreateSrv(srvc)
					m[srv.GetAddr()] = srv
					go srv.Run()
					time.Sleep(100 * time.Millisecond)
					if !srv.CheckOk() {
						delete(m, srv.GetAddr())
					}
				}
				//fmt.Printf("%v\n", *srv)
			} else if ev.Type == db.DELOP {
				fmt.Println("in delete op")
				go srv.Stop()
				delete(m, res[0][2])
			}
		}
	}
}
