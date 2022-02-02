package qcache

import (
	"fmt"
	"log"
	"sync"
)

// ===================================== 流程图 =====================================
//							  是
//	接收 key --> 检查是否被缓存 -----> 返回缓存值[1]
//					| 否                         是
//	                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值[2]
//				    |  否
//					|-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值[3]
// ===================================== 流程图 =====================================

// Group
// 	缓存命名空间
// 	数据加载
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("Getter should not be <nil>")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name] // 找不到，则为 <nil>
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if val, ok := g.mainCache.get(key); ok {
		log.Println("[QCache]", key, "hit")
		return val, nil
	}
	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// ===================================== 流程图 =====================================
//	使用一致性哈希选择节点            是                                    是
//  		|-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
//						|  否                                    ↓  否
//  					|--------------------------------> 回退到本地节点处理
// ===================================== 流程图 =====================================
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[QCache] Failed to get from peer", err)
		}
	}
	// 本地获取
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.set(key, value)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// Getter 从数据源加载数据（不在缓存中）
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数[接口型函数]
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
