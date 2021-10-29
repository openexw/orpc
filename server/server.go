package server

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/openexw/orpc/internal/protocol"
	"io"
	"log"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	ErrServerClosed = errors.New("rpc#server: Server closed")
)

const (
	ReaderBufSize = 1024
)

type Option struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	trace        bool // 判断是否开启 trace
}

type OptionFn func(opts *Option)

type OptionFns []OptionFn

func (fns OptionFns) apply(opt *Option) {
	for _, fn := range fns {
		fn(opt)
	}
}

// WithReadTimeout can set read timeout
// read timeout default is 3 second
func WithReadTimeout(timeout time.Duration) OptionFn {
	return func(opts *Option) {
		opts.readTimeout = timeout
	}
}

// WithWriteTimeout can set write timeout
// write timeout default is 3 second
func WithWriteTimeout(timeout time.Duration) OptionFn {
	return func(opts *Option) {
		opts.writeTimeout = timeout
	}
}

// WithTrace If you want open debug mode, It can do this
// trace default is closed
func WithTrace(isTrace bool) OptionFn {
	return func(opts *Option) {
		opts.trace = isTrace
	}
}

// defaultOption default option
var defaultOption = Option{
	readTimeout:  300 * time.Second,
	writeTimeout: 300 * time.Second,
	trace:        false,
}

// Server 是整个 oprc 的核心功能
type Server struct {
	*Option
	ln net.Listener

	serviceMap sync.Map //map[string]*service
	mu         sync.RWMutex
}

// NewServer 实例化 server 对象
//
//	s := NewServer(
//		WithReadTimeout(time.Second),	// 设置读超时
//		WithWriteTimeout(time.Second),	// 设置写超时
//		WithTrace(true))	// 开启 trace
func NewServer(opt ...OptionFn) *Server {
	options := defaultOption
	OptionFns(opt).apply(&options)

	return &Server{Option: &options}
}

// Server 可以启动一个服务端，提供了一个 Dail() 的方法启动
//
// 	s := NewServer()
// 	err := s.Server("tcp", ":9081")
func (server *Server) Server(network, address string) (err error) {
	var ln net.Listener
	ln, err = net.Listen(network, address)
	if err != nil {
		return
	}
	log.Println("start rpc server on", ln.Addr())
	return server.serveListener(ln)
}

// serveListener 监听客户端连接
func (server *Server) serveListener(ln net.Listener) error {
	server.mu.Lock()
	server.ln = ln
	server.mu.Unlock()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// 针对连接关闭做出特殊处理
			if strings.Contains(err.Error(), "listener closed") {
				return ErrServerClosed
			}
			return err
		}
		go server.serveConn(conn)
	}
}

// serveConn 处理连接
func (server *Server) serveConn(conn net.Conn) {
	r := bufio.NewReaderSize(conn, ReaderBufSize)
	// TODO ..
	wg := sync.WaitGroup{}
	for {
		t := time.Now()
		// 设置 timeout
		if server.readTimeout != 0 {
			_ = conn.SetReadDeadline(t.Add(server.readTimeout))
		}

		// 获取一个请求
		request, err := server.readRequest(r)
		if server.trace {
			log.Printf("rpc##server: request:%v err: %v\n", request, err)
		}
		if err != nil {
			if err == io.EOF {
				log.Printf("rpc##server: client has closed this connection: %s\n", conn.RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("rpc##server: client has closed this connection: %s\n", conn.RemoteAddr().String())
			} else {
				log.Printf("rpc##server: client has closed this connection: %s\n", conn.RemoteAddr().String())
			}
			return
		}

		// 设置写超时
		if server.writeTimeout != 0 {
			_ = conn.SetWriteDeadline(t.Add(server.writeTimeout))
		}
		wg.Add(1)
		//处理请求
		go func() {
			// 捕获异常
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					err := fmt.Errorf("[rpc##server handler request error]: %v, request: %v, stack: %s",
						r, request, buf)
					log.Println(err)
				}
			}()
			// 处理请求
			resp, err := server.handleRequest(request)
			if server.trace {
				log.Printf("rpc##server: handler request after, resp:%v => err:%v from conn: %s\n", resp, err, conn.RemoteAddr().String())
			}
			if err != nil {
				log.Printf("rpc##server failed to handle request: %v\n", err)
				return
			}
			// 响应
			rst := resp.Encode()
			_, err = conn.Write(rst)
			if server.trace {
				log.Println("rpc##server: send data to conn")
			}
			if err != nil {
				log.Println("rpc: server connect write err:", err)
			}
		}()
	}
}

// readRequest 获取一个请求
func (server *Server) readRequest(r io.Reader) (*protocol.Message, error) {
	req := protocol.NewMessage()
	err := req.Decode(r)
	if server.trace {
		log.Printf("readRequest header %v, ServiceMethod :%s, Payload:%s. error:%v", req.Header, req.ServiceMethod, req.Payload, err)
	}
	// 获取到到一个请求
	return req, err
}

// handleRequest 处理从客户端发送的请求，并通过反射机制调用 service.method
func (server *Server) handleRequest(request *protocol.Message) (resp *protocol.Message, err error) {
	srv, mType, err := server.findService(request.ServiceMethod)
	if err != nil {
		return nil, err
	}

	var newReq *protocol.Message
	newReq = request
	newReq.SetMessageType(protocol.Response)

	// 序列化
	codec := protocol.Codecs[request.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", request.SerializeType())
		return
	}

	// 处理参数
	arg := mType.newArg()
	argv := arg.Interface()
	// 获取请求参数
	err = codec.Decode(request.Payload, &argv)
	if err != nil {
		return
	}
	if server.trace {
		log.Printf("rpc##server: get a service %+v for an request %+v", srv, argv)
	}

	// 处理返回值
	reply := mType.newReply()

	// 调用函数 service.call
	if mType.ArgType.Kind() != reflect.Ptr {
		err = srv.call(mType, arg.Elem(), reply)
	} else {
		err = srv.call(mType, arg, reply)
	}
	replyv := reply.Interface()
	if server.trace {
		log.Printf("rpc##server: reply = %v", replyv)
	}
	if err != nil {
		return nil, err
	}
	if replyv != nil {
		data, err := codec.Encode(replyv)
		if err != nil {
			return nil, err
		}
		newReq.Payload = data
	}
	return newReq, nil
}

// findService 根据 serviceMethod  从 service.serviceMethod 中获取 service 信息
func (server *Server) findService(serviceMethod string) (srv *service, mType *methodType, err error) {
	index := strings.LastIndex(serviceMethod, ".")
	if index < 0 {
		err = errors.New("rpc##server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]
	svri, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc#server: can't find service " + serviceName)
		return
	}

	srv = svri.(*service)
	if mType = srv.method[methodName]; mType == nil {
		err = errors.New("rpc##server: can't find method " + serviceMethod)
		return
	}
	return
}

// Register service 中满足条件的方法集合，Method 必须满足以下条件才能被注册：
//
// 	1. 方法可导出
// 	2. 必须包含两个参数，一个参数、一个返回值，并且可被导出
// 	3. 参数和返回值都必须是指针
// 	4. 必须包含一个返回值且是 error 类型
func (server *Server) Register(rcvr interface{}) error {
	s, err := newService(rcvr)
	if err != nil {
		return err
	}
	if _, loaded := server.serviceMap.LoadOrStore(s.name, s); loaded {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

// Close 关闭 server 连接
func (server *Server) Close() error {
	server.mu.Lock()
	defer server.mu.Unlock()

	var err error

	if server.ln != nil {
		err = server.ln.Close()
	}
	return err
}

// Dail 是快捷启动服务方法，也可以通过以下方式启动：
//
//	s := NewServer(WithTrace(true))
//	s.Server("tcp", ":9089")
func Dail(network, address string, opts ...OptionFn) error {
	server := NewServer(opts...)
	return server.Server(network, address)
}
