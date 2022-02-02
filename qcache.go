package QCache

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

func (g *Group) load(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.mainCache.set(key, value)
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
