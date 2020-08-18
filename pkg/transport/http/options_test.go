package http

import (
	"strconv"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/specs"
)

func TestParseEndpointOptions(t *testing.T) {
	duration := time.Second
	method := "POST"
	endpoint := "/endpoint"
	codec := "xml"

	options := specs.Options{
		MethodOption:        method,
		EndpointOption:      endpoint,
		ReadTimeoutOption:   duration.String(),
		WriteTimeoutOption:  duration.String(),
		RequestCodecOption:  codec,
		ResponseCodecOption: codec,
	}

	result, err := ParseEndpointOptions(options)
	if err != nil {
		t.Fatal(err)
	}

	if result.Method != method {
		t.Errorf("unexpected method %+v, expected %+v", result.Method, method)
	}

	if result.Endpoint != endpoint {
		t.Errorf("unexpected endpoint %+v, expected %+v", result.Endpoint, endpoint)
	}

	if result.ReadTimeout != duration {
		t.Errorf("unexpected read timeout %+v, expected %+v", result.ReadTimeout, duration)
	}

	if result.WriteTimeout != duration {
		t.Errorf("unexpected write timeout %+v, expected %+v", result.ReadTimeout, duration)
	}

	if result.RequestCodec != codec {
		t.Errorf("unexepected request codec %+v, expected %+v", result.RequestCodec, codec)
	}

	if result.ResponseCodec != codec {
		t.Errorf("unexepected response codec %+v, expected %+v", result.RequestCodec, codec)
	}
}

func TestParseEndpointOptionsRequestResponseCodec(t *testing.T) {
	codec := "xml"

	options := specs.Options{
		CodecOption: codec,
	}

	result, err := ParseEndpointOptions(options)
	if err != nil {
		t.Fatal(err)
	}

	if result.RequestCodec != codec {
		t.Errorf("unexepected request codec %+v, expected %+v", result.RequestCodec, codec)
	}

	if result.ResponseCodec != codec {
		t.Errorf("unexepected response codec %+v, expected %+v", result.RequestCodec, codec)
	}
}

func TestParseCallerOptions(t *testing.T) {
	duration := time.Second
	idle := 500

	options := specs.Options{
		MaxIdleConnsOption:  strconv.Itoa(idle),
		TimeoutOption:       duration.String(),
		KeepAliveOption:     duration.String(),
		FlushIntervalOption: duration.String(),
	}

	result, err := ParseCallerOptions(options)
	if err != nil {
		t.Fatal(err)
	}

	if result.MaxIdleConns != idle {
		t.Fatalf("unexpected max idle connections %+v, expected %+v", result.MaxIdleConns, idle)
	}

	if result.Timeout != duration {
		t.Fatalf("unexpected timeout %+v, expected %+v", result.Timeout, duration)
	}

	if result.KeepAlive != duration {
		t.Fatalf("unexpected keep alive %+v, expected %+v", result.KeepAlive, duration)
	}

	if result.FlushInterval != duration {
		t.Fatalf("unexpected flush interval %+v, expected %+v", result.FlushInterval, duration)
	}
}
