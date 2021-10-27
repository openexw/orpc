package testdata

import "github.com/openexw/orpc/internal/protocol"

var (
	Str           = []byte("Hello world")
	ServiceMethod = "Foo.Sum"
)

func BuildMessageData() []byte {
	m := protocol.NewMessage()
	m.ServiceMethod = ServiceMethod
	m.SetVersion(1)
	m.SetMessageType(protocol.Request)
	m.SetSeq(uint64(1))
	m.SetSerializeType(protocol.Msgpack)
	m.CheckMagicNumber()
	m.SetCompressType(protocol.Gzip)
	m.Payload = Str
	return m.Encode()
}
