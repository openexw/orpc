package protocol

import (
	"github.com/openexw/orpc/testdata"
	"testing"
)

func TestRegisterCodec(t *testing.T) {
	mockCodec := testdata.NewMockCodec()
	mockSerializeType := 100

	err := RegisterCodec(SerializeType(mockSerializeType), mockCodec)
	if err != nil {
		t.Errorf("RegisterCodec error:%v", err)
		return
	}

	//assert.e
}
