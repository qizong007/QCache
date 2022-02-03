# 简介

一个**http缓存库**

（本项目仅作为学习使用）

# 特性

- 以`lru`为缓存结构作淘汰
- 支持`并发`，并可解决`缓存击穿`问题
- 支持`分布式`
    - 加入了`一致性哈希`
    - 使用`HTTP`通信
    - 以`protobuf`作为通信格式
- `HTTP`方式使用

# 项目参考

- GeeCache
- group-cache（Go版的`Memcached`）