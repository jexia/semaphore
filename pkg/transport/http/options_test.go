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

	options := specs.Options{
		MethodOption:       method,
		EndpointOption:     endpoint,
		ReadTimeoutOption:  duration.String(),
		WriteTimeoutOption: duration.String(),
	}

	result, err := ParseEndpointOptions(options)
	if err != nil {
		t.Fatal(err)
	}

	if result.Method != method {
		t.Fatalf("unexpected method %+v, expected %+v", result.Method, method)
	}

	if result.Endpoint != endpoint {
		t.Fatalf("unexpected endpoint %+v, expected %+v", result.Endpoint, endpoint)
	}

	if result.ReadTimeout != duration {
		t.Fatalf("unexpected read timeout %+v, expected %+v", result.ReadTimeout, duration)
	}

	if result.WriteTimeout != duration {
		t.Fatalf("unexpected write timeout %+v, expected %+v", result.ReadTimeout, duration)
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
