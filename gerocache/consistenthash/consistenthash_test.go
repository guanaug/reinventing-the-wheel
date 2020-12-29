package consistenthash

import (
	"log"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	fn := func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}
	hash := New(3, fn)

	hash.Add("6", "4", "2")

	keys := map[string]string{
		"2": "2", "11": "2", "23": "4", "27": "2",
	}

	for k, v := range keys {
		if hash.Get(k) != v {
			log.Fatalf("should get %v, but get %v", v, hash.Get(k))
		}
	}

	hash.Add("8")
	keys["27"] = "8"
	for k, v := range keys {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
