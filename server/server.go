package server

//import (
//	"github.com/openexw/orpc/codec"
//	"github.com/openexw/orpc/protocol"
//	"io"
//	"log"
//	"net"
//	"sync"
//)
//
//const (
//	// Seq 默认值
//	Seq            = 0x3bef5c
//	ReaderBuffSize = 1024
//)
//
//type options struct {
//	Seq           int                    // 标识本次请求的唯一序列号
//	SerializeType protocol.SerializeType // 请求的类型
//}
//type Server struct {
//	lock sync.Mutex
//	wg   sync.WaitGroup
//	options
//}
//
//// 默认的选项配置
//var defaultOptions = options{
//	Seq:           Seq,
//	SerializeType: protocol.Gob,
//}
//
//type OptFunc func(*options)
//
//type OptsFunc []OptFunc
//
//func (fns OptsFunc) apply(opt *options) {
//	for _, fn := range fns {
//		fn(opt)
//	}
//}
//
//// WithSerializeType 配置序列化类型
//func WithSerializeType(t protocol.SerializeType) OptFunc {
//	return func(options *options) {
//		options.SerializeType = t
//	}
//}
//
//// WithSeq 配置唯一序列号
//func WithSeq(seq int) OptFunc {
//	return func(options *options) {
//		options.Seq = seq
//	}
//}
//
//// NewServer 实例化一个 server
//// 可使用 `NewServer(WithSeq(1), WithSerializeType())` 配置参数实例化 server
//func NewServer(opts ...OptFunc) *Server {
//	options := defaultOptions
//	OptsFunc(opts).apply(&options)
//
//	return &Server{options: options}
//}
//
//func (s *Server) Accept(listen net.Listener) {
//	for {
//		// 阻塞等待 socket 连接
//		conn, err := listen.Accept()
//		if err != nil {
//			log.Println("orpc server: accept error:", err)
//			return
//		}
//
//		// 处理连接
//		go s.ServeConn(conn)
//	}
//}
//
//// Serve 处理连接成功的请求
//func (s *Server) ServeConn(conn net.Conn) {
//	//r := bufio.NewReaderSize(conn, ReaderBuffSize)
//
//	//
//}
//
//// serveCodec 处理请求，分为三个阶段
//// 1. 读请求
//// 2. 处理请求
//// 3. 发送响应
//func (s *Server) serveCodec(c codec.Codec) {
//	for {
//		req, err := s.readRequest(c)
//		if err != nil {
//			if req == nil {
//				break // it's not possible to recover, so close the connection
//			}
//			req.Error = err.Error()
//			s.sendResponse(c, req.h, invalidRequest, sending)
//			continue
//		}
//		s.wg.Add(1)
//		go s.handleRequest(c, req, sending, wg)
//	}
//	s.wg.Wait()
//	_ = c.Close()
//}
//
//// readRequestHeader 读取请求头
//func (s *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
//	var ch *codec.Header
//	// 读取请求头
//	if err := c.ReaderHeader(ch); err != nil {
//		if err != io.EOF && err != io.ErrUnexpectedEOF {
//			log.Println("orpc server: read header error:", err)
//		}
//		return nil, err
//	}
//
//	return ch, nil
//}
//
//// readRequest 读取请求
//func (s *Server) readRequest(c codec.Codec) (*codec.Header, error) {
//	header, err := s.readRequestHeader(c)
//	if err != nil {
//		return nil, err
//	}
//
//}
//
//func (s *Server) sendResponse(c interface{}, h interface{}, request interface{}, sending interface{}) {
//
//}
//
//func (s *Server) handleRequest(c interface{}, req interface{}, sending interface{}, wg interface{}) {
//
//}
