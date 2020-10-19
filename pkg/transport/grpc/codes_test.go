package grpc

import (
	"net/http"
	"strconv"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestCodeFromStatus(t *testing.T) {
	tests := map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           http.StatusRequestTimeout,
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.Unauthenticated:    http.StatusUnauthorized,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
	}

	for code, expected := range tests {
		t.Run(strconv.Itoa(int(expected)), func(t *testing.T) {
			result := StatusFromCode(code)
			if result != expected {
				t.Fatalf("unexepcted result %+v, expected %+v", result, expected)
			}
		})
	}
}

func TestStatusFromCode(t *testing.T) {
	tests := map[int]codes.Code{
		http.StatusOK:                  codes.OK,
		http.StatusRequestTimeout:      codes.Canceled,
		http.StatusBadRequest:          codes.InvalidArgument,
		http.StatusGatewayTimeout:      codes.DeadlineExceeded,
		http.StatusNotFound:            codes.NotFound,
		http.StatusConflict:            codes.AlreadyExists,
		http.StatusForbidden:           codes.PermissionDenied,
		http.StatusUnauthorized:        codes.Unauthenticated,
		http.StatusTooManyRequests:     codes.ResourceExhausted,
		http.StatusNotImplemented:      codes.Unimplemented,
		http.StatusInternalServerError: codes.Internal,
		http.StatusServiceUnavailable:  codes.Unavailable,
	}

	for code, expected := range tests {
		t.Run(strconv.Itoa(int(code)), func(t *testing.T) {
			result := CodeFromStatus(code)
			if result != expected {
				t.Fatalf("unexepcted result %+v, expected %+v", result, expected)
			}
		})
	}
}
