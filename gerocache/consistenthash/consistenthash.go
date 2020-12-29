package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash struct {
	fn       func(data []byte) uint32
	replicas int
	keys     []int          // hash 环
	hash     map[int]string // key->host
}

func New(replicas int, fn func(data []byte) uint32) *Hash {
	if nil == fn {
		fn = crc32.ChecksumIEEE
	}

	return &Hash{
		fn:       fn,
		replicas: replicas,
		hash:     make(map[int]string),
	}
}

func (h *Hash) Add(hosts ...string) {
	for _, host := range hosts {
		for i := 0; i < h.replicas; i++ {
			key := int(h.fn([]byte(strconv.Itoa(i) + host)))
			h.keys = append(h.keys, key)
			h.hash[key] = host
		}
	}

	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})
}

func (h *Hash) Get(key string) string {
	if 0 == len(h.keys) {
		return ""
	}

	hash := int(h.fn([]byte(key)))
	idx := sort.SearchInts(h.keys, hash)

	return h.hash[h.keys[idx%len(h.keys)]]
}
