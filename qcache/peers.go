package qcache

// PeerPicker 用于根据传入的 key 选择相应节点的PeerGetter（HTTP服务端）
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 用于从对应 group 查找缓存值（HTTP客户端）
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
