package server

import (
	"bufio"
	"bytes"
	"github.com/openexw/orpc/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

var s *Server

func init() {
	s = NewServer()
}
func TestServer_readRequest(t *testing.T) {
	md := testdata.BuildMessageData()
	buf := bytes.NewBuffer(md)
	r := bufio.NewReaderSize(buf, 1024)
	request, err := s.readRequest(r)
	if err != nil {
		t.Error("read request error:", err)
	}
	if request == nil {
		t.Error("request is nil")
	}
}

func TestServer_Register(t *testing.T) {
	s := NewServer()
	s.Register(new(testdata.Sum))

	var counter int
	s.serviceMap.Range(func(key, value interface{}) bool {
		counter++
		return true
	})
	assert.Equal(t, 1, counter)
}
