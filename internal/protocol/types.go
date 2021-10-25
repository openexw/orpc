package protocol

import (
	"errors"
	"github.com/openexw/orpc/internal/codec"
)

type SerializeType byte

//GobType  SerializeType = "application/gob"
//JsonType SerializeType = "application/json"
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
	Protobuf
)

// Codecs 是当前支持的 codec 列表，如有需要可以通过 RegisterCodec(t SerializeType, cc codec.Codec) 注册自定义的 Codecs
var Codecs = map[SerializeType]codec.Codec{
	Raw:     codec.NewRawCodec(),
	Msgpack: codec.NewMsgpackCodec(),
	Gob:     codec.NewGobCodec(),
	JSON:    codec.NewJsonCodec(),
	//Protobuf:
}

var codecExist = errors.New("SerializeType is exist")

// RegisterCodec 注册自定义的 codec
func RegisterCodec(t SerializeType, cc codec.Codec) error {
	if _, ok := Codecs[t]; ok {
		return codecExist
	}
	Codecs[t] = cc
	return nil
}
