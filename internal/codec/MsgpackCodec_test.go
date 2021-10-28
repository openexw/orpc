package codec

import (
	"github.com/openexw/orpc/testdata"
	"testing"
)

func TestMsgpackCodec_Encode(t *testing.T) {
	m := NewMsgpackCodec()
	encode, err := m.Encode(testdata.JsonProfileStr)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("%s", encode)
	//type args struct {
	//	i interface{}
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	want    []byte
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		m := MsgpackCodec{}
	//		got, err := m.Encode(tt.args.i)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("Encode() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
