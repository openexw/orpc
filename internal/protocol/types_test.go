package protocol

import (
	"github.com/openexw/orpc/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterCodec(t *testing.T) {
	mockCodec := testdata.NewMockCodec()
	mockSerializeType := 100
	originCodecLen := len(Codecs)

	err := RegisterCodec(SerializeType(mockSerializeType), mockCodec)
	if err != nil {
		t.Errorf("RegisterCodec error:%v", err)
		return
	}

	assert.Equal(t, originCodecLen+1, len(Codecs))
}

func TestRegisterCodec_Error(t *testing.T) {
	mockCodec := testdata.NewMockCodec()
	mockSerializeType := 3
	originCodecLen := len(Codecs)
	err := RegisterCodec(SerializeType(mockSerializeType), mockCodec)

	assert.Error(t, errCodecExist, err)
	assert.Equal(t, originCodecLen, len(Codecs))
}

func TestMagicNumber(t *testing.T) {
	tests := []struct {
		name string
		want byte
	}{
		{
			name: "test_get_magicNumber",
			want: magicNumber,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MagicNumber(); got != tt.want {
				t.Errorf("MagicNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
