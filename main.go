package main

import (
	"fmt"
	"net"
	"regexp"

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

	srvs := make(map[string]server.CommonServer)

	for _, srvc := range allSrvConf {
		// 防止重复的key加入到map中
		srv, err := server.CreateSrv(srvc)
		if err != nil {
			fmt.Println("create server error:", err)
			continue
		}
		_, ok := srvs[srv.GetAddr()]
		if ok {
			fmt.Printf("%s server is  already existed", srv.GetAddr())
			continue
		}
		srvs[srv.GetAddr()] = srv
		go srv.Run()
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

			re := regexp.MustCompile(`among/server/(\w+)/(.+)`)
			res := re.FindAllStringSubmatch(string(ev.Kv.Key), -1)
			if len(res[0]) < 3 {
				continue
			}
			srv := m[res[0][2]]

			if ev.Type == db.PUTOP {
				if ev.PrevKv != nil { // modify existed one
					if res[0][1] == server.MySQLSTR {
						mc := new(server.MySQLServerConf)
						err := mc.Unmarshal(ev.Kv.Value)
						if err != nil {
							fmt.Println("put op config error", err)
							continue
						}
						if srv == nil {
							fmt.Println("in put new one op, srv is nil", res[0][2])
							srvType := res[0][1]
							createSrvRun(srvType, ev.Kv.Value, m)
							fmt.Printf("m is %v\n", m)
							continue
						}

						if srv.CompareConf(mc) { //这里还有一个问题是，以前是mysql srv然后端口没改，变成了redis srv，其实应该把srv delete了重新生成
							fmt.Println("among db not changed")
							continue
						}

						srv.SetSrvPort(mc.SrvPort)
					}
					go srv.Reload()
				} else { // add new one
					srvType := res[0][1]
					createSrvRun(srvType, ev.Kv.Value, m)
				}
			} else if ev.Type == db.DELOP {
				if srv == nil {
					fmt.Println("in del op, srv is nil")
					continue
				}
				go srv.Stop()
				delete(m, res[0][2])
			}
		}
	}
}

func createSrvRun(ty string, conf []byte, m map[string]server.CommonServer) {
	srvc := server.CreateSrvConf(ty, conf)
	if srvc == nil {
		fmt.Println("add new server conf error")
		return
	}
	srv, err := server.CreateSrv(srvc)
	if err != nil {
		return
	}
	m[srv.GetAddr()] = srv
	go srv.Run()
}
