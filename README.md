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
|     1byte   |  1byte  |    1byte     |      1byte   |     1byte     | 8byte |
|-------------|---------|--------------|--------------|---------------|-------|
```

说明：

- MagicNumber：用于校验
- Version：版本号
- MessageType：消息类型，分为 `Request` 和 `Response`
- CompressType：压缩类型，可为 `None` 和 `Gzip`
- SerializeType：序列化类型，现支持 `Raw`、`JSON`、`Msgpack`、`Gob`
- Seq：序列号，存储的是一个 `32` 位的无符号整型

## 通信

通信原理图如图所示：

![y5rqLboRPnUuKTv](https://i.loli.net/2021/10/28/y5rqLboRPnUuKTv.png)

## Feature list

- RPC 协议定义，协议包含 Header 和 Message ，详细见 「消息体」定义
- Codec 序列化协议支持，默认为 Msgpack，支持 Gob、JSON、Msgpack、Protobuf（暂未实现）
- 支持自定义 Codec，可通过 `protocol.RegisterCodec()` 注册自定义的 codec
- 支持设置压缩（暂未实现，已经设计了该功能）
- 实现 Client、Server 端
- 目前仅支持 TCP 和 UDP 协议

## TODO

- [ ] feat：负载均衡功能
- [ ] feate：支持 HTTP、HTTPS、WS 等协议
- [ ] feat：分布式链路最终功能
- [ ] feat：服务发现模块
- [ ] feat：新增插件机制
- [ ] optimize：sync.Pool 优化
- [ ] optimize：rpc 协议中的 header 空间浪费问题，部分可使用 `1bit` 存储

## 快速开始

> **Notice**：
>
> 1. 详细的 example 可查看 `./example` 中的 demo。
> 2. 请在 Go 1.7 下运行 
> **运行服务端**
>
> ```shel
> $ go run example/server.go
> ```
>
> **运行客户端**
>
> ```shell
> $ go run example/client.go
> ```

**server.go**

```go
package main

import (
	"github.com/openexw/orpc/server"
	"github.com/openexw/orpc/testdata"
	"log"
)

func main() {
	// 实例化 server
	s := server.NewServer()
	// 注册 service
	s.Register(new(testdata.Sum))
	// 启动 server
	err := s.Server("tcp", ":8091")
	if err != nil {
		log.Fatalln("run orpc err:", err)
	}
}
```

**testdata/sum.go**

```go
package testdata

type Sum int
type Args struct {
	A int
	B int
}

func (s *Sum) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}
```

**client.go**

```go
package main

import (
	"context"
	"fmt"
	"github.com/openexw/orpc/client"
	"github.com/openexw/orpc/testdata"
	"net"
)

// Args 定义请求参数
type Args struct {
	A, B int
}

func main() {
	conn, err := net.Dial("tcp", ":8091")
	if err != nil {
		return
	}

	// 实例化客户端
	cli := client.NewClient(conn, client.WithIsTrace(false))

	// 设置请求参数
	args := &Args{
		A: 12,
		B: 2,
	}
	// 返回值
	var reply int
	// call
	cli.Call(context.Background(), "Sum.Add", args, &reply)
	// 输出结果
	println(args.A, "+", args.B, "=", reply)
}
```

## 参考

- net/rpc
- rpcx
