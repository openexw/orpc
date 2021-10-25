package orpc

import (
	"net"
	"testing"
)

func TestServer_Accept(t *testing.T) {
	listen, _ := net.Listen("tcp", ":9081")
	s := NewServer()
	s.Accept(listen)
}
