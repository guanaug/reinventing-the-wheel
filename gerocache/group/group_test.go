package group

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup(t *testing.T) {
	loadCount := make(map[string]int, len(db))
	group := NewGroup("scores", 1<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search for key", key)
		if v, ok := db[key]; ok {
			loadCount[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Fatalf("should be get value %s, but get %s", v, view.String())
		}
		if view, err := group.Get(k); err != nil || loadCount[k] > 1 || view.String() != v {
			log.Fatalf("cache %s miss", k)
		}
	}

	if view, err := group.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
