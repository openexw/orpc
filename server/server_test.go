package server

import (
	"bufio"
	"bytes"
	"github.com/openexw/orpc/internal/protocol"
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
func TestServer_batchReadRequest(t *testing.T) {
	for i := 0; i < 10; i++ {
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
}

func TestServer_Register(t *testing.T) {
	s.Register(new(testdata.Sum))

	var counter int
	s.serviceMap.Range(func(key, value interface{}) bool {
		counter++
		return true
	})
	assert.Equal(t, 1, counter)
}

func TestServer_handleRequest(t *testing.T) {
	s.Register(new(testdata.Profile))
	req := testdata.BuildRequest()
	resp, err := s.handleRequest(req)
	if err != nil {
		t.Fatalf("failed to hand request: %v", err)
	}

	if resp.Payload == nil {
		t.Fatalf("expect reply but got %s", resp.Payload)
	}
	var profile testdata.Profile
	codec := protocol.Codecs[resp.SerializeType()]
	if codec == nil {
		t.Fatalf("can not find codec %c", codec)
	}

	err = codec.Decode(resp.Payload, &profile)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	assert.Equal(t, uint8(18), profile.Age)
	assert.Equal(t, uint8(1), profile.Sex)
}

func TestServer_findService(t *testing.T) {
	_, _, err := s.findService("Foo.Sum")
	if err == nil {
		assert.Equal(t, nil, err)
	}

	err = s.Register(new(testdata.Profile))
	if err != nil {
		t.Fatalf("Can't register service : %v", err)
	}
	srv, _, err := s.findService("Profile.AddProfile")
	if err != nil {
		t.Fatalf("Can't find service : %v", err)
	}

	if srv == nil {
		t.Fatalf("service is nill")
	}
}
