- 架构图
  - 功能列表
    - 负载均衡
      - 轮询
      - 加权轮询
      - 随机
    - RPC 协议
      - 组成
        - header [13]byte
          - MagicNumber
          - Version
          - MessageType
          - SerializeType
          - CompressType 
          - Seq
        - Payload
      - 传输协议
        - TCP
        - HTTP / HTTPS
        - UDP
      - 序列化协议
        - protobuf
        - json
        - gob
        - Msgpack
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

## 消息体
### Message 协议
```
|--------|---------|------------------|--------------------|------------|---------------|
| header | dataLen | serviceMethodLen |    serviceMethod   | payloadLen |    payload    |
|--------|---------|------------------|--------------------|------------|---------------|
| 13byte |  4byte  |       4byte      | len(serviceMethod) |    4byte   |  len(payload) |
|--------|---------|------------------|--------------------|------------|---------------|
```
### Header 格式
```
|-------------|---------|--------------|--------------|---------------|-------|
| MagicNumber | Version | MessageType  | CompressType | SerializeType |  Seq  |
|-------------|---------|--------------|--------------|---------------|-------|
|     1byte   |  1byte  |    1byte     |      1byte   |       1byte   | 8byte |
|-------------|---------|--------------|--------------|---------------|-------|
```
说明：
- MagicNumber：用于校验
- Version：版本号
- MessageType：消息类型，分为 Request 和 Response
- CompressType：压缩类型，可为 Note 和 Gzip
- SerializeType：序列化类型，现支持 Raw、JSON、Msgpack、Gob
- Seq：序列号，存储的是一个 32 位的无符号整型

## 通信
## 启动
```go
s := orpc.NewServer()
s.Register("Foo.sum")
s.Server("tcp", ":8990")
```