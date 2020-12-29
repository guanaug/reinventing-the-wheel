package communication

import (
	"fmt"
	"gerocache/consistenthash"
	"gerocache/group"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultPrefix   = "/_gerocache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	addr  string
	hash  *consistenthash.Hash
	peers map[string]*httpGetter
	mu    sync.Mutex
}

func NewHTTPPool(baseURL string) *HTTPPool {
	return &HTTPPool{
		addr:  baseURL,
		hash:  consistenthash.New(defaultReplicas, nil),
		peers: make(map[string]*httpGetter),
	}
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, defaultPrefix) {
		http.Error(w, "unrecognized url", http.StatusBadRequest)
		return
	}

	// http://[ip:port]/$defaultPrefix/$group/$key
	urls := strings.SplitN(r.URL.Path[len(defaultPrefix):], "/", 2)
	if len(urls) < 2 {
		http.Error(w, "group or key is missing", http.StatusBadRequest)
		return
	}
	grpName, key := urls[0], urls[1]

	grp := group.GetGroup(grpName)
	if nil == grp {
		http.Error(w, "group not found", http.StatusBadRequest)
		return
	}

	view, err := grp.Get(key)
	if err != nil {
		log.Printf("get key error: %v, group is: %v, key is: %s\n", err, grp, key)
		http.Error(w, "get key error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(view.ByteSlice())
}

func (h *HTTPPool) RegisterPeer(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.hash.Add(peers...)
	for _, peer := range peers {
		h.peers[peer] = &httpGetter{
			baseURL: peer + defaultPrefix,
		}
	}
}

// TODO
func (h *HTTPPool) PickPeer(key string) (group.PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if peer := h.hash.Get(key); peer != "" && peer != h.addr {
		log.Println("pickup:", peer)
		return h.peers[peer], true
	}

	return nil, false
}

type httpGetter struct {
	baseURL string
}

func (h httpGetter) Get(group string, key string) ([]byte, error) {
	query := fmt.Sprintf("http://%s%s/%s", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	resp, err := http.Get(query)
	if err != nil {
		log.Printf("Get error, url: %s", query)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("server return code %d", resp.StatusCode)
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
