package main

import (
	"QCache/qcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	qcache.NewGroup("scores", 2<<10, qcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:7777"
	peers := qcache.NewHTTPPool(addr)
	log.Println("qcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}


