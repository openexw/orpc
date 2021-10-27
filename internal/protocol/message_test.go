package protocol

import (
	"bytes"
	"github.com/openexw/orpc/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	encode := testdata.BuildMessageData()
	buf := bytes.NewBuffer(encode)

	m := NewMessage()
	err := m.Decode(buf)
	if err != nil {
		t.Errorf("Decode err: %v", err)
		return
	}
	assert.Equal(t, testdata.Str, m.Payload)
	assert.Equal(t, testdata.ServiceMethod, m.ServiceMethod)
}
