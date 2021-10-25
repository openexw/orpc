package codec

import (
	"bytes"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackCodec struct {
}

func NewMsgpackCodec() *MsgpackCodec {
	return &MsgpackCodec{}
}

// Encode 将对象编码成 []byte
func (m MsgpackCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(msgpack.Marshaler); ok {
		return m.MarshalMsgpack()
	}

	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

// Decode 将 []byte 解析成对象
func (m MsgpackCodec) Decode(data []byte, i interface{}) error {
	if m, ok := i.(msgpack.Unmarshaler); ok {
		return m.UnmarshalMsgpack(data)
	}
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	return dec.Decode(i)
}
