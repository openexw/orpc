package protocol

import "io"

//type Header [12]byte

const (
	magicNumber byte = 0x08
)

type Header struct {
	//MagicNumber, Version, MessageType, SerializeType, MessageStatusType, CompressType, Heartbeat, Oneway byte
	MagicNumber byte
}
type Message struct {
	*Header
	ServiceMethod string //格式形如：Service.Method
	Payload       []byte
	data          []byte
}

func NewMessage(serviceMethod string) *Message {
	return &Message{Header: &Header{MagicNumber: magicNumber}}
}

func (m *Message) Encode() []byte {
	// TODO ...
	return nil
}
func (m *Message) Decode(r io.Reader) error {
	// TODO ...
	return nil
}
