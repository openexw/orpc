package codec

import (
	"bytes"
	"encoding/gob"
)

// GobCodec 使用 encoding/gob 进行编码
type GobCodec struct {
}

func NewGobCodec() *GobCodec {
	return &GobCodec{}
}

// Encode 将对象编码成 []byte
func (g GobCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

// Decode 将对象解码
func (g GobCodec) Decode(data []byte, i interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	return decoder.Decode(i)
}
