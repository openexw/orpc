package testdata

type MockCodec struct {

}

func NewMockCodec() *MockCodec {
	return &MockCodec{}
}

func (m MockCodec) Encode(i interface{}) ([]byte, error) {
	return nil,nil
}

func (m MockCodec) Decode(data []byte, i interface{}) error {
	return nil
}



