package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 一致性哈希
type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 虚拟节点倍数 => 解决数据倾斜
	keys     []int          // 哈希环（虚拟节点编号）
	hashMap  map[int]string // 虚拟节点 -> 真实节点
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // default
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 生成 replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// hash值作为虚拟节点编号
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys) // 保持有序
}

// Get 获取下一个最接近的节点（32位hash -> 虚拟节点）
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 找下一个最接近的节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 如果idx == len(m.keys)，应选择 m.keys[0]，因为这是哈希"环"
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
