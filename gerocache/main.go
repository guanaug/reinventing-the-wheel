package main

import (
	"flag"
	"fmt"
	"gerocache/communication"
	"gerocache/group"
	"log"
	"net/http"
	"strings"
)

func init() {
	log.SetPrefix("[GeroCache] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var db = map[string]string{
	"k1": "v1", "k2": "v2", "k3": "v3",
}

func startCacheServer(addr string, peers ...string) error {
	fn := group.GetterFunc(func(key string) ([]byte, error) {
		if val, ok := db[key]; ok {
			fmt.Printf("[SlowDB] search key: %s, get val: %s", key, val)
			return []byte(val), nil
		}
		return nil, fmt.Errorf("key not found")
	})

	// TODO 不能动态增删 peers
	pool := communication.NewHTTPPool(addr)
	pool.RegisterPeer(peers...)
	group.NewGroupPeers("group", 1<<10, fn, pool)

	fmt.Println("start server on:", addr)
	return http.ListenAndServe(addr, pool)
}

func usage() {
	fmt.Println(`./gerocache -addr ip:port -peers ip1:port1,ip2:port2`)
}

func main() {
	addr := flag.String("addr", "", "-addr 127.0.0.1:8000")
	peer := flag.String("peers", "", "-peers 127.0.0.1:8000,127.0.0.1:8001")
	flag.Parse()

	if "" == *addr || "" == *peer || len(strings.Split(*peer, ",")) < 1 {
		usage()
		return
	}

	peers := strings.Split(*peer, ",")
	peers = append(peers, *addr)

	log.Fatal(startCacheServer(*addr, peers...))
}
