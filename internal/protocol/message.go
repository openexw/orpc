package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

//type iHeader interface {
//	// CheckMagicNumber 校验 magicNumber，使用 1byte 存储
//	CheckMagicNumber() bool
//
//	// Version 获取版本号，使用 1byte 存储
//	Version() byte
//	// SetVersion 设置版本号
//	SetVersion(version byte)
//
//	// MessageType 获取消息类型，使用 1byte 存储
//	MessageType() MessageType
//	// SetMessageType 设置消息类型
//	SetMessageType(messageType MessageType)
//
//	// compressType 获取压缩类型，使用 1byte 存储
//	compressType() compressType
//	// SetCompressType 设置压缩类型
//	SetCompressType(compressType compressType)
//
//	// serializeType 获取序列化类型，使用 1byte 存储
//	serializeType() serializeType
//	// SetSerializeType 设置序列化类型
//	SetSerializeType(serializeType serializeType)
//}

var ErrMagicNumber = errors.New("magic number error")

// Header 是 Message 的头部信息
// Format：
// |-------------|---------|--------------|--------------|---------------|-------|
// | MagicNumber | Version | MessageType  | CompressType | SerializeType |  Seq  |
// |-------------|---------|--------------|--------------|---------------|-------|
// |     1byte   |  1byte  |    1byte     |      1byte   |       1byte   | 8byte |
// |-------------|---------|--------------|--------------|---------------|-------|
// MagicNumber：用于校验
// Version：版本号
// MessageType：消息类型，分为 Request 和 Response
// CompressType：压缩类型，可为 Note 和 Gzip
// SerializeType：序列化类型，现支持 Raw、JSON、Msgpack、Gob
// Seq：序列号，存储的是一个 8 位的整数
type Header [HeaderLen]byte

// SerializeType 序列化类型，默认支持 Raw、Gob、JSON、Msgpack
type SerializeType byte

// MessageType 包括 Request（请求）和 Response（响应）两种类型
type MessageType byte

// CompressType 包括 None（不压缩）和 Gzip 压缩类型
type CompressType byte

// CheckMagicNumber 校验 magicNumber，使用 1byte 存储
func (h *Header) CheckMagicNumber() bool {
	return h[0] == magicNumber
}

// Version 获取版本号，使用 1byte 存储
func (h *Header) Version() byte {
	return h[1]
}

// SetVersion 设置版本号
func (h *Header) SetVersion(version byte) {
	h[1] = version
}

// MessageType 获取消息类型，使用 1byte 存储
func (h *Header) MessageType() MessageType {
	return MessageType(h[2])
}

// SetMessageType 设置消息类型
func (h *Header) SetMessageType(messageType MessageType) {
	h[2] = byte(messageType)
}

// CompressType 获取压缩类型，使用 1byte 存储
func (h *Header) CompressType() CompressType {
	return CompressType(h[3])
}

// SetCompressType 设置压缩类型
func (h *Header) SetCompressType(compressType CompressType) {
	h[3] = byte(compressType)
}

// SerializeType 获取序列化类型，使用 1byte 存储
func (h *Header) SerializeType() SerializeType {
	return SerializeType(h[4])
}

// SetSerializeType 设置序列化类型
func (h *Header) SetSerializeType(serializeType SerializeType) {
	h[4] = byte(serializeType)
}

// Seq 获取序列号
func (h *Header) Seq() uint64 {
	return binary.BigEndian.Uint64(h[5:])
}

// SetSeq 设置序列号
func (h *Header) SetSeq(seq uint64) {
	binary.BigEndian.PutUint64(h[5:], seq)
}

// Message 传输的消息格式
// format:
// |--------|---------|------------------|--------------------|------------|---------------|
// | header | dataLen | serviceMethodLen |    serviceMethod   | payloadLen |    payload    |
// |--------|---------|------------------|--------------------|------------|---------------|
// | 13byte |  4byte  |       4byte      | len(serviceMethod) |    4byte   |  len(payload) |
// |--------|---------|------------------|--------------------|------------|---------------|
//
// data = serviceMethodLen + serviceMethod + payloadLen + payload
type Message struct {
	*Header
	ServiceMethod string //格式形如：Service.Method
	Payload       []byte // 真实的内容载体
	data          []byte
}

// TODO 错误处理

// NewMessage 实例化一个 message 对象
func NewMessage() *Message {
	header := Header([HeaderLen]byte{})
	header[0] = magicNumber
	return &Message{
		Header: &header,
	}
}

func (m *Message) Encode() []byte {
	smLen := len(m.ServiceMethod) // "Service.Method" 的长度
	// 			header + dataLen + (sml+ServiceMethod) + payloadLen + payload
	dataLen := ServiceMethodLen + smLen + PayloadLen + len(m.Payload)
	data := make([]byte, 4096)
	// step1: 写入 header
	copy(data, m.Header[:])

	// step2：写入 DataLen
	binary.BigEndian.PutUint32(data[HeaderLen:HeaderLen+DataLen], uint32(dataLen))

	// step3：写入 sml
	// 写入 ServiceMethod 的长度
	serviceMethodLenStart := HeaderLen + DataLen
	binary.BigEndian.PutUint32(data[serviceMethodLenStart:serviceMethodLenStart+ServiceMethodLen], uint32(smLen))
	// 写入 ServiceMethod 的值
	serviceMethodStart := serviceMethodLenStart + ServiceMethodLen
	copy(data[serviceMethodStart:serviceMethodStart+smLen], m.ServiceMethod)

	//step4：写入 payload
	payloadLen := len(m.Payload)
	payloadStart := serviceMethodStart + smLen
	// 写入长度
	binary.BigEndian.PutUint32(data[payloadStart:payloadStart+PayloadLen], uint32(payloadLen))
	// 写入数据
	copy(data[payloadStart+PayloadLen:], m.Payload)
	return data
}

func (m *Message) Decode(r io.Reader) error {
	// step1：解析 header
	_, err := io.ReadFull(r, m.Header[:])
	if err != nil {
		return err
	}
	// 校验 MagicNumber
	if !m.Header.CheckMagicNumber() {
		//log.Printf("rpc##server message.Decode err:%s", ErrMagicNumber.Error()+"; value:"+string(m.Header[0]))
		return ErrMagicNumber
	}
	// step2: 获取 data
	// 获取 dataLen
	dataLenByte := make([]byte, DataLen)
	_, err = io.ReadFull(r, dataLenByte)
	if err != nil {
		return err
	}
	dataLen := binary.BigEndian.Uint32(dataLenByte)
	dataTotalLen := int(dataLen)
	// 获取 data
	if cap(m.data) >= dataTotalLen {
		m.data = m.data[:dataTotalLen]
	} else {
		m.data = make([]byte, dataTotalLen)
	}

	// step3：获取 ServiceMethod
	data := m.data
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}
	// ServiceMethod length
	smLenByte := data[0:ServiceMethodLen]
	smLen := int(binary.BigEndian.Uint32(smLenByte))
	// ServiceMethod
	m.ServiceMethod = string(data[ServiceMethodLen : ServiceMethodLen+smLen])

	// step4: 获取 Payload
	payloadStart := ServiceMethodLen + smLen
	// 获取 payloadLen
	//payloadLenByte := data[payloadStart:PayloadLen]
	//payloadLen := int(binary.BigEndian.Uint32(payloadLenByte))

	// 获取 Payload
	payload := data[payloadStart+PayloadLen:]
	m.Payload = payload
	return nil
}
