package group

import (
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name       string
	localCache *cache
	getter     Getter
	peers      PeerPick
}

var (
	groups = make(map[string]*Group)
	grpMu  = sync.RWMutex{}
)

func NewGroup(name string, maxBytes int64, getter Getter) *Group {
	return NewGroupPeers(name, maxBytes, getter, nil)
}

func NewGroupPeers(name string, maxBytes int64, getter Getter, peers PeerPick) *Group {
	if nil == getter {
		panic("Getter must be not nil")
	}

	grpMu.Lock()
	defer grpMu.Unlock()

	groups[name] = &Group{
		name:       name,
		localCache: newCache(maxBytes),
		getter:     getter,
		peers:      peers,
	}

	return groups[name]
}

func (g *Group) RegisterPeers(peers PeerPick) {
	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	// 本地缓存找到
	if val, ok := g.localCache.cache.Get(key); ok {
		return val.(ByteView), nil
	}

	// 本地缓存没找到，且配置了分布式节点，根据一致性hash看在哪个节点找
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if bv, err := g.getFromRemote(peer, key); nil == err {
				return bv, err
			}
			log.Printf("get key %s, from peer: %v error", key, peer)
		}
	}

	// 远程获取失败，则尝试从数据源获取
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	val, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	view := ByteView{data: val}
	g.populateCache(key, view)
	return view, nil
}

func (g *Group) getFromRemote(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{bytes}, nil
}

func (g *Group) populateCache(key string, view ByteView) {
	g.localCache.set(key, view)
}

func GetGroup(name string) *Group {
	grpMu.RLock()
	defer grpMu.RUnlock()
	return groups[name]
}
