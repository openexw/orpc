package codec

import (
	"github.com/openexw/orpc/codec/testdata"
	"reflect"
	"testing"
)

func TestJsonCodec_Encode(t *testing.T) {
	j := NewJsonCodec()
	type Args struct {
		i interface{}
	}
	tests := []struct {
		Name    string
		Args    Args
		Want    []byte
		WantErr bool
	}{
		{
			Name:    "single_profile_encode",
			Args:    Args{i: testdata.ProfileSingle},
			Want:    []byte(testdata.JsonProfileStr),
			WantErr: false,
		},
		{
			Name:    "empty_profile_encode",
			Args:    Args{i: testdata.Profile{}},
			Want:    []byte(testdata.JsonZeroStr),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := j.Encode(tt.Args.i)
			if (err != nil) != tt.WantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.WantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.Want)
			}
		})
	}
}

func TestJsonCodec_Decode(t *testing.T) {
	j := NewJsonCodec()
	type args struct {
		data []byte
		i    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "single_profile_decode",
			args: args{
				data: []byte(testdata.JsonProfileStr),
				i:    &testdata.Profile{},
			},
			wantErr: false,
		},
		{
			name: "empty_profile_decode",
			args: args{
				data: []byte(testdata.JsonZeroStr),
				i:    &testdata.Profile{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := j.Decode(tt.args.data, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v, %+v", err, tt.wantErr, tt.args.i)
			}
		})
	}
}

// BenchmarkJsonCodec_Encode JsonCodec encode 基准测试
func BenchmarkJsonCodec_Encode(b *testing.B) {
	var raw = make([]byte, 0, 1024)
	j := NewJsonCodec()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw, _ = j.Encode(testdata.ProfileSingle)
	}
	b.ReportMetric(float64(len(raw)), "bytes")
}

// BenchmarkJsonCodec_Decode JsonCodec decode 基准测试
func BenchmarkJsonCodec_Decode(b *testing.B) {
	j := NewJsonCodec()

	var rst testdata.Profile
	str := []byte(testdata.JsonProfileStr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = j.Decode(str, &rst)
	}
}
