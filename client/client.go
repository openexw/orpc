package client

import (
	"bufio"
	"context"
	"errors"
	"github.com/openexw/orpc/internal/protocol"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrShutdown         = errors.New("connection is shut down")
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

// ReaderBufferSize 客户端读取数据的 buf size
const ReaderBufferSize = 16 * 1024

// Options 是客户端的连接选项
type Options struct {
	// connTimeout 表示连接超时
	connTimeout time.Duration
	// compressType 表示压缩类型
	compressType protocol.CompressType
	// serializeType 表示序列化类型
	serializeType protocol.SerializeType
	// isTrace 表示是否是调试模式
	isTrace bool
}

// defaultOptions 默认的配置
var defaultOptions = Options{
	// 默认超时时间为 3 s
	connTimeout: 3 * time.Second,
	// 默认不进行压缩
	compressType: protocol.None,
	// 默认使用 Msgpack 进行序列化
	serializeType: protocol.Msgpack,
	isTrace:       false,
}

type OptionFn func(options *Options)

type OptionFns []OptionFn

func (fns OptionFns) apply(opt *Options) {
	for _, fn := range fns {
		fn(opt)
	}
}

// WithConnTimeout 设置超时时间
func WithConnTimeout(timeout time.Duration) OptionFn {
	return func(opt *Options) {
		opt.connTimeout = timeout
	}
}

// WithCompressType 设置压缩类型
func WithCompressType(compressType protocol.CompressType) OptionFn {
	return func(opt *Options) {
		opt.compressType = compressType
	}
}

// WithSerializeType 设置序列化类型
func WithSerializeType(serializeType protocol.SerializeType) OptionFn {
	return func(opt *Options) {
		opt.serializeType = serializeType
	}
}

func WithIsTrace(isTrace bool) OptionFn {
	return func(opt *Options) {
		opt.isTrace = isTrace
	}
}

type Call struct {
	ServiceMethod string
	Args          interface{}
	Reply         interface{}

	Done  chan *Call
	Error error
	Ctx   context.Context
}

// Client ORPC 连接客户端
type Client struct {
	*Options          // 客户端连接选项
	conn     net.Conn // 连接

	reqMutex sync.Mutex // 发送 req 保护
	mu       sync.Mutex
	seq      uint64           // 序列号
	pending  map[uint64]*Call // 调用列表

	closing  bool      // 通过 close 正常不关闭客户端
	shutdown bool      // 异常退出，比如报错
	r        io.Reader // 读取请求的 reader buf
}

// NewClient 实例化
// example：
//
// client := NewClient(
//		WithCompressType(protocol.Gzip),
//		WithSerializeType(protocol.JSON))
func NewClient(conn net.Conn, opts ...OptionFn) *Client {
	options := defaultOptions
	OptionFns(opts).apply(&options)

	c := &Client{
		Options: &options,
		r:       bufio.NewReaderSize(conn, ReaderBufferSize),
		pending: make(map[uint64]*Call),
		conn:    conn,
	}

	go c.dealResponse()
	return c
}

// dealResponse 处理从 server 端来的响应
func (client *Client) dealResponse() {
	var err error
	for err == nil {
		// 响应
		resp := protocol.NewMessage()
		err = resp.Decode(client.r)
		if err != nil {
			break
		}

		// 处理非服务端请求
		seq := resp.Seq()

		isServerMessage := resp.MessageType() == protocol.Request
		var call *Call
		if !isServerMessage {
			call = client.remove(seq)
		}

		// TODO deal with message err
		switch {
		case call == nil:
			if isServerMessage {
				err = errors.New("call is nil")
				continue
			}
		default:
			data := resp.Payload
			if len(data) > 0 {
				// 对响应的数据进行解码（反序列化）
				codec := protocol.Codecs[resp.SerializeType()]
				if codec == nil {
					call.Error = ErrUnsupportedCodec
				} else {
					// 反序列化到 call.Reply
					if err = codec.Decode(data, call.Reply); err != nil {
						call.Error = err
					}
				}
			}
			call.done()
		}
	}

	// 客户端关闭，终止操作
	client.terminate(err)
}

// terminate 客户端关闭，终止操作
func (client *Client) terminate(err error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	closing := client.closing
	client.conn.Close()
	client.shutdown = true
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}

	for _, call := range client.pending {
		call.Error = err
		call.done()
	}

	if err != nil && !closing {
		log.Println("orpc: client protocol error:", err)
	}
}

// register add call to pending
func (client *Client) register(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	// 检测客户端的状态
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	return seq, nil
}

// remove delete call from pending
func (client *Client) remove(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()

	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

// Call a rpc service
func (client *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) (err error) {
	call := client.Go(ctx, serviceMethod, args, reply)

	// 处理 Error
	select {
	case <-ctx.Done():
		rmCall := client.remove(client.seq)
		if rmCall != nil {
			err = ctx.Err()
			call.done()
		}
	case callDone := <-call.Done:
		err = callDone.Error
	}
	return err
}

// Go 会初始化一个 Call 并发送 Call 到 server
func (client *Client) Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) *Call {
	// net/rpc 为啥会初始化两次
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          make(chan *Call, 1),
		Ctx:           ctx,
	}
	// send data to conn
	client.send(call)
	return call
}

// send 发送数据到 server
func (client *Client) send(call *Call) {
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()

	// step1: 注册  到pending
	seq, err := client.register(call)
	if err != nil {
		call.Error = err
		call.done()
	}

	// step2：处理请求的 header 和 message 参数
	req := protocol.NewMessage()
	req.SetSeq(seq)
	req.SetMessageType(protocol.Request)
	req.SetSerializeType(client.Options.serializeType)
	// TODO 加上压缩
	//req.SetCompressType()
	req.ServiceMethod = call.ServiceMethod

	// step3：处理请求参数（Payload）
	// 获取 codec
	codec := protocol.Codecs[client.Options.serializeType]
	if codec == nil {
		call.Error = ErrUnsupportedCodec
		call.done()
		return
	}
	// 序列化参数
	data, err := codec.Encode(call.Args)
	if err != nil {
		rmCall := client.remove(seq)
		if rmCall != nil {
			call.Error = err
			call.done()
		}
		return
	}
	req.Payload = data
	// 编码参数
	reqData := req.Encode()
	if client.isTrace {
		log.Printf("rpc##client req: header is %v, ServiceMethod is %s, paylod is %s,  client send data\n", req.Header, req.ServiceMethod, req.Payload)
	}
	// step4: 发送数据
	_, err = client.conn.Write(reqData)
	if err != nil {
		rmCall := client.remove(seq)
		if rmCall != nil {
			call.Error = err
			call.done()
		}
		return
	}
}

// done call done
func (call *Call) done() {
	select {
	case call.Done <- call:
		//	ok
	default:
		log.Println("rpc##client: discarding Call reply due to insufficient Done chan capacity")
	}
}
