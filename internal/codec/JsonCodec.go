package codec

import (
	"bytes"
	"encoding/json"
)

// JsonCodec 使用 Marshal / NewDecoder 进行编码
type JsonCodec struct {
}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

// Encode 将对象编码成一个 []byte
func (j JsonCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

// Decode 从一个 []byte 解码成一个对象
func (j JsonCodec) Decode(data []byte, i interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	return decoder.Decode(i)
}
