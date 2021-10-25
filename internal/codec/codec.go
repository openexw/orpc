package codec

// Codec 定义协议 decode/encode payload
type Codec interface {
	// Encode 编码
	Encode(i interface{}) ([]byte, error)
	// Decode 解码
	Decode(data []byte, i interface{}) error
}
