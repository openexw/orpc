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

func BuildRequest() *protocol.Message {
	m := protocol.NewMessage()
	m.ServiceMethod = "Profile.AddProfile"
	m.SetMessageType(protocol.Request)
	m.SetSerializeType(protocol.Msgpack)
	m.SetSeq(uint64(12))

	arg := &Profile{
		Name: "Jack",
		Age:  18,
		Sex:  1,
	}

	codec := protocol.Codecs[protocol.Msgpack]
	data, _ := codec.Encode(arg)

	m.Payload = data
	return m
}
