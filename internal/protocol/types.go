package protocol

import (
	"errors"
	"github.com/openexw/orpc/internal/codec"
)

const (
	// HeaderLen Message header length
	HeaderLen = 13
	// DataLen Message DataLen length
	DataLen = 4
	// PayloadLen Message Payload length
	PayloadLen = 4
	// ServiceMethodLen Message ServiceMethod length
	ServiceMethodLen = 4
)

// SerializeType 序列化类型，默认支持 Raw、Gob、JSON、Msgpack
const (
	// Raw 不进行序列化，使用 []byte
	Raw SerializeType = iota
	// Gob 使用 encode/gob 进行序列化
	Gob
	// JSON 使用 encode/json 进行序列化
	JSON
	// Msgpack 使用 Msgpack 进行序列化
	Msgpack
	// Protobuf 使用 ProtoBuf 进行序列化
	//Protobuf
)

// Codecs 是当前支持的 codec 列表，如有需要可以通过 RegisterCodec(t SerializeType, cc codec.Codec) 注册自定义的 Codecs
var Codecs = map[SerializeType]codec.Codec{
	Raw:     codec.NewRawCodec(),
	Gob:     codec.NewGobCodec(),
	JSON:    codec.NewJsonCodec(),
	Msgpack: codec.NewMsgpackCodec(),
	//Protobuf:
}

var errCodecExist = errors.New("SerializeType is exist")

// RegisterCodec 注册自定义的 codec
func RegisterCodec(t SerializeType, cc codec.Codec) error {
	if _, ok := Codecs[t]; ok {
		return errCodecExist
	}
	Codecs[t] = cc
	return nil
}

const (
	// magicNumber 用于请求的校验
	magicNumber byte = 0x13
)

// MagicNumber 获取 MagicNumber
func MagicNumber() byte {
	return magicNumber
}

// MessageType 包括 Request（请求）和 Response（响应）两种类型
const (
	// Request 请求
	Request MessageType = iota

	// Response 响应
	Response
)

// CompressType 包括 None（不压缩）和 Gzip 压缩类型
const (
	// None 不压缩
	None CompressType = iota
	// Gzip 使用 gzip 进行压缩
	Gzip
)
