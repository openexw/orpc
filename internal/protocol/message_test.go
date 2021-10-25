package protocol

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	str           = []byte("Hello world")
	serviceMethod = "Foo.Sum"
)

func buildMessageData() []byte {

	m := NewMessage()
	m.ServiceMethod = serviceMethod
	m.SetVersion(1)
	m.SetMessageType(Request)
	m.SetSeq(uint64(1))
	m.SetSerializeType(Msgpack)
	m.CheckMagicNumber()
	m.SetCompressType(Gzip)
	m.Payload = str
	return m.Encode()
}

func TestHeader_header(t *testing.T) {
	mType := Response
	version := byte(100)
	sType := Msgpack
	cType := Gzip
	seq := uint64(1)

	m := NewMessage()
	m.SetVersion(version)
	m.SetMessageType(mType)
	m.SetSerializeType(sType)
	m.SetCompressType(cType)
	m.SetSeq(seq)

	assert.Equal(t, mType, m.MessageType())
	assert.Equal(t, version, m.Version())
	assert.Equal(t, cType, m.CompressType())
	assert.Equal(t, sType, m.SerializeType())
	assert.Equal(t, seq, m.Seq())
	assert.Equal(t, true, m.CheckMagicNumber())
}

func TestMessage_Decode(t *testing.T) {
	encode := buildMessageData()
	buf := bytes.NewBuffer(encode)

	m := NewMessage()
	err := m.Decode(buf)
	if err != nil {
		t.Errorf("Decode err: %v", err)
		return
	}
	assert.Equal(t, str, m.Payload)
	assert.Equal(t, serviceMethod, m.ServiceMethod)
}
