package grpc

import (
	"fmt"
	"testing"
)

func ValidateBytes(expected []byte, result []byte) error {
	if len(expected) != len(result) {
		return fmt.Errorf("unexpected result %+v, expected %+v", result, expected)
	}

	for index := range result {
		if expected[index] != result[index] {
			return fmt.Errorf("unexpected byte %+v, expected %+v in %+v : %+v", result[index], expected[index], result, expected)
		}
	}

	return nil
}

func TestCodec(t *testing.T) {
	framePtr := &frame{}
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	codec := rawCodec{}

	err := codec.Unmarshal(data, framePtr)
	if err != nil {
		t.Fatal(err)
	}

	out, err := codec.Marshal(framePtr)
	if err != nil {
		t.Fatal(err)
	}

	err = ValidateBytes(data, out)
	if err != nil {
		t.Fatal(err)
	}

	in := []byte{0x55}
	err = codec.Unmarshal(in, framePtr)
	if err != nil {
		t.Fatal(err)
	}

	out, err = codec.Marshal(framePtr)
	if err != nil {
		t.Fatal(err)
	}

	err = ValidateBytes(in, out)
	if err != nil {
		t.Fatal(err)
	}
}
