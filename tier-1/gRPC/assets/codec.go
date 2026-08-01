package assets

import "fmt"

// RawCodec passes []byte payloads through untouched, skipping Protobuf
// marshaling so the gRPC benchmark measures pure binary transport.
type RawCodec struct{}

func (RawCodec) Marshal(v any) ([]byte, error) {
	b, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("RawCodec: message is %T, want []byte", v)
	}
	return b, nil
}

func (RawCodec) Unmarshal(data []byte, v any) error {
	b, ok := v.(*[]byte)
	if !ok {
		return fmt.Errorf("RawCodec: message is %T, want *[]byte", v)
	}
	*b = data
	return nil
}

func (RawCodec) Name() string {
	return "raw"
}
