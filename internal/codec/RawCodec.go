package codec

import (
	"fmt"
	"reflect"
)

// RawCodec 使用 []byte 作为原始格式
type RawCodec struct {
}

func NewRawCodec() *RawCodec {
	return &RawCodec{}
}

// Encode 返回 []byte 数组
func (r RawCodec) Encode(i interface{}) ([]byte, error) {
	// []type 类型
	if data, ok := i.([]byte); ok {
		return data, nil
	}

	// 为指针类型
	if data, ok := i.(*[]byte); ok {
		return *data, nil
	}
	return nil, fmt.Errorf("%T is not a []byte", i)
}

// Decode 返回 []byte 数组
func (r RawCodec) Decode(data []byte, i interface{}) error {
	// 获取 i 中的值，并设置给 data
	reflect.Indirect(reflect.ValueOf(i)).SetBytes(data)
	return nil
}
