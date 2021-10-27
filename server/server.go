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
	"strings"
	"sync"
	"time"
)

var (
	ErrServerClosed = errors.New("http: Server closed")
)

// TODO ADD options add isTrace

type Server struct {
	ln           net.Listener
	readTimeout  time.Duration
	writeTimeout time.Duration

	serviceMap sync.Map //map[string]*service
	//mu sync.RWMutex
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Server(network, address string) (err error) {
	var ln net.Listener
	ln, err = net.Listen(network, address)
	if err != nil {
		return
	}
	log.Println("start rpc server on", ln.Addr())
	return server.serveListener(ln)
}

func (server *Server) serveListener(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// 针对连接关闭做出特殊处理
			if strings.Contains(err.Error(), "listener closed") {
				return ErrServerClosed
			}
			return err
		}
		// TODO add go
		go server.serveConn(conn)
	}
}

func (server *Server) serveConn(conn net.Conn) {
	log.Println("rpc##server serveConn")
	r := bufio.NewReaderSize(conn, 1024)

	for {
		t := time.Now()
		// 设置 timeout
		if server.readTimeout != 0 {
			_ = conn.SetReadDeadline(t.Add(server.readTimeout))
		}

		// 获取一个请求
		request, err := server.readRequest(r)
		if request == nil {
			//log.Println("rpc##server: Nil request")
			break
		}
		log.Printf("rpc##server: request:%v err: %v", request, err)
		if err != nil {
			// TODO
			if err == io.EOF {
				//log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				//log.Infof("rpcx: connection %s is closed", conn.RemoteAddr().String())
			} else {
				//log.Warnf("rpcx: failed to read request: %v", err)
			}
			continue
		}

		// 设置写超时
		if server.writeTimeout != 0 {
			_ = conn.SetWriteDeadline(t.Add(server.writeTimeout))
		}

		// 处理请求
		go func() {
			//defer func() {
			//	if r := recover(); r != nil {
			//	}
			//}()
			log.Println("rpc: server handler request before")
			// 处理请求
			resp, err := server.handleRequest(request)
			log.Printf("rpc: server handler request after, resp:%v => err:%v\n", resp, err)
			if err != nil {
				log.Printf("rpc: server failed to handle request: %v\n", err)
			}
			// 响应
			rst := resp.Encode()
			_, err = conn.Write(rst)
			log.Printf("write data to conn\n")
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
	// 获取到到一个请求
	if len(req.Payload) == 0 {
		return nil, err
	}
	return req, err
}

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
	log.Printf("handleRequest: %v", argv)

	// 处理返回值
	reply := mType.newReply()

	// 调用函数 service.call
	if mType.ArgType.Kind() != reflect.Ptr {
		err = srv.call(mType, arg.Elem(), reply)
	} else {
		err = srv.call(mType, arg, reply)
	}
	replyv := reply.Interface()
	log.Printf("rpc##server: reply = %v", replyv)
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
	}

	srv = svri.(*service)
	if mType = srv.method[methodName]; mType == nil {
		err = errors.New("rpc##server: can't find method " + serviceMethod)
		return
	}
	return
}

// Register service 中满足条件的方法集合，方法有以下条件：
// Method 必须满足一下条件：
// 1. 方法可导出
// 2. 必须包含两个参数，一个参数、一个返回值，并且可被导出
// 3. 参数和返回值都必须是指针
// 4. 必须包含一个返回值且是 error 类型
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, loaded := server.serviceMap.LoadOrStore(s.name, s); loaded {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}
