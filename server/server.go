package server

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/openexw/orpc/internal/protocol"
	"io"
	"log"
	"net"
	"strings"
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

	//serviceMap []
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
	log.Println("rpc: server serveConn")
	// TODO
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
			continue
		}
		log.Printf("rpc##server: request:%v err: %v", request, err)
		//if err != nil {
		//	// TODO
		//	if err == io.EOF {
		//		//log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
		//	} else if strings.Contains(err.Error(), "use of closed network connection") {
		//		//log.Infof("rpcx: connection %s is closed", conn.RemoteAddr().String())
		//	} else {
		//		//log.Warnf("rpcx: failed to read request: %v", err)
		//	}
		//	continue
		//}
		//
		//// 设置写超时
		//if server.writeTimeout != 0 {
		//	_ = conn.SetWriteDeadline(t.Add(server.writeTimeout))
		//}
		//
		//// 处理请求
		//go func() {
		//	//defer func() {
		//	//	if r := recover(); r != nil {
		//	//	}
		//	//}()
		//	log.Println("rpc: server handler request before")
		//	// 处理请求
		//	resp, err := server.handleRequest(request)
		//	log.Printf("rpc: server handler request after, resp:%v => err:%v\n", resp, err)
		//	if err != nil {
		//		log.Printf("rpc: server failed to handle request: %v\n", err)
		//	}
		//	// 响应
		//	rst := resp.Encode()
		//	_, err = conn.Write(rst)
		//	log.Printf("write data to conn\n")
		//	if err != nil {
		//		log.Println("rpc: server connect write err:", err)
		//	}
		//}()
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
	serviceMethod := request.ServiceMethod
	index := strings.LastIndex(serviceMethod, ".")
	if index < 0 {
		err = errors.New("rpc##server: service/method request ill-formed: " + serviceMethod)
		return
	}
	//serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]

	var newReq *protocol.Message
	newReq = request
	newReq.SetMessageType(protocol.Response)

	// TODO deal serviceName and methodName

	// 序列化
	codec := protocol.Codecs[request.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", request.SerializeType())
		return
	}

	// 获取请求参数
	// TODO ERROR deal
	//var argv reflect.Type
	//err = codec.Decode(request.Payload, &argv)
	//if err != nil {
	//	return
	//}
	//log.Printf("handleRequest: %v", argv)
	replv := fmt.Sprintf("orpc resp %d", request.Seq())
	data, err := codec.Encode([]byte(replv))
	if err != nil {
		return nil, err
	}
	newReq.Payload = data // TODO ..
	return newReq, nil
}
