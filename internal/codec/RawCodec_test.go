package codec

import (
	"encoding/json"
	"github.com/openexw/orpc/testdata"
	"reflect"
	"testing"
)

func TestRawCodec_Encode(t *testing.T) {
	r := NewRawCodec()
	pBytes, _ := json.Marshal(testdata.ProfileSingle)
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "single_profile_encode",
			args:    args{i: pBytes},
			want:    []byte(testdata.JsonProfileStr),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Encode(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawCodec_Decode(t *testing.T) {
	r := NewRawCodec()
	rst := new([]byte)

	_ = r.Decode(*rst, testdata.JsonProfileStr)

	t.Logf("%v", rst)
	//
	//type args struct {
	//	data []byte
	//	i    interface{}
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	wantErr bool
	//}{
	//	{
	//		name: "single_profile_raw_decode",
	//		args: args{
	//			data: rst,
	//			i:    testdata.JsonProfileStr,
	//		},
	//		wantErr: false,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if err := r.Decode(tt.args.data, tt.args.i); (err != nil) != tt.wantErr {
	//			t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
	//		}
	//	})
	//}
}
