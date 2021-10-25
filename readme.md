- 架构图
  - 功能列表
    - 负载均衡
      - 轮询
      - 加权轮询
      - 随机
    - RPC 协议
      - 组成
        - header [5]byte
          - MagicNumber
          - Version
          - MessageType
          - SerializeType
          - CompressType 
        - body
      - 传输协议
        - TCP
        - HTTP / HTTPS
        - UDP
      - 序列化协议
        - protobuf
        - json
        - gob
        - MsgPack
        - Raw -> 原始格式 []byte
        - thrift
    - 监控
    - 分布式链路追踪
      - prometheus
    - 服务发现
      - 技术选型
        - zookeeper
        - etcd
        - consul
- 通信图
- 
## 架构图
![Apm4KtBSHG8Jx53](https://i.loli.net/2021/10/22/Apm4KtBSHG8Jx53.png)


## 设计
- 使用 `encoding/gob` 实现消息的编码（序列化和反序列化操作）

## 通信
一次 TCP 通信，至少需要三部分类容，Options、Header、Body，结构如下：
![TUDuqpa2GoONIxP](https://i.loli.net/2021/10/20/TUDuqpa2GoONIxP.png)

- options 采用固定格式 json 存储，标识本次请求的唯一序列号、编码类型（god、json等）
- header 和 body 格式如上交替发送，可以发送多个 header 和 body 的组合

options 的结构体如下：

```go
type Options struct {
	Seq       int        // 标识本次请求的唯一序列号
	CodecType codec.Type // 请求的类型
}
```

## 启动
```go
s := orpc.NewServer()
s.Register("User", new(User), "")
s.Server("tcp", ":8990")
```