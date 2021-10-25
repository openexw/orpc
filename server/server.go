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
	"time"
)

var (
	ErrServerClosed = errors.New("http: Server closed")
)

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

		server.serveConn(conn)
	}

}

func (server *Server) serveConn(conn net.Conn) {
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
		if err != nil {
			// TODO
			if err == io.EOF {
				//log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				//log.Infof("rpcx: connection %s is closed", conn.RemoteAddr().String())
			} else {
				//log.Warnf("rpcx: failed to read request: %v", err)
			}
			return
		}

		// 设置写超时
		if server.writeTimeout != 0 {
			_ = conn.SetWriteDeadline(t.Add(server.writeTimeout))
		}

		// 处理请求
		go func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			// 处理请求
			resp, err := server.handleRequest(request)
			if err != nil {
				log.Printf("rpcx: failed to handle request: %v\n", err)
			}
			// 响应
			rst := resp.Encode()
			conn.Write(rst)
		}()
	}
}

// readRequest 获取一个请求
func (server *Server) readRequest(r *bufio.Reader) (req *protocol.Message, err error) {
	req = protocol.NewMessage()
	err = req.Decode(r)
	// 获取到到一个请求
	if err == io.EOF {
		return req, err
	}
	return nil, err
}

func (server *Server) handleRequest(request *protocol.Message) (resp *protocol.Message, err error) {
	serviceMethod := request.ServiceMethod
	index := strings.LastIndex(serviceMethod, ".")
	if index < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
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
	var argv reflect.Type
	err = codec.Decode(request.Payload, &argv)
	if err != nil {
		return
	}
	log.Printf("handleRequest: %v", argv)
	replv := fmt.Sprintf("orpc resp %d", request.Seq())
	newReq.Payload = []byte(replv) // TODO ..
	return newReq, nil
}
