package formencoded

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jexia/semaphore/pkg/specs/types"
)

// AddTypeKey encodes the given value into the given encoder
func AddTypeKey(encoder url.Values, key string, typed types.Type, value interface{}) {
	var encoded string
	switch typed {
	case types.Double:
		encoded = Float64Empty(value)
	case types.Int64:
		encoded = Int64Empty(value)
	case types.Uint64:
		encoded = Uint64Empty(value)
	case types.Fixed64:
		encoded = Uint64Empty(value)
	case types.Int32:
		encoded = Int32Empty(value)
	case types.Uint32:
		encoded = Uint32Empty(value)
	case types.Fixed32:
		encoded = Uint32Empty(value)
	case types.Float:
		encoded = Float32Empty(value)
	case types.String:
		encoded = StringEmpty(value)
	case types.Enum:
		encoded = StringEmpty(value)
	case types.Bool:
		encoded = BoolEmpty(value)
	case types.Bytes:
		encoded = BytesBase64Empty(value)
	case types.Sfixed32:
		encoded = Int32Empty(value)
	case types.Sfixed64:
		encoded = Int64Empty(value)
	case types.Sint32:
		encoded = Int32Empty(value)
	case types.Sint64:
		encoded = Int64Empty(value)
	}

	encoder.Add(key, encoded)
}

// DecodeType decodes the given property from the given decoder
func DecodeType(decoder url.Values, key string, typed types.Type) (interface{}, error) {
	raw := decoder.Get(key)

	switch typed {
	case types.Double:
		value, err := strconv.ParseFloat(raw, 64)
		return value, err
	case types.Float:
		value, err := strconv.ParseFloat(raw, 64)
		return float64(value), err
	case types.Int64:
		value, err := strconv.ParseInt(raw, 10, 64)
		return value, err
	case types.Uint64:
		value, err := strconv.ParseUint(raw, 10, 64)
		return value, err
	case types.Fixed64:
		value, err := strconv.ParseUint(raw, 10, 64)
		return value, err
	case types.Int32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case types.Uint32:
		value, err := strconv.ParseUint(raw, 10, 32)
		return uint32(value), err
	case types.Fixed32:
		value, err := strconv.ParseUint(raw, 10, 32)
		return uint32(value), err
	case types.String:
		return raw, nil
	case types.Bool:
		return strconv.ParseBool(raw)
	case types.Bytes:
		value := make([]byte, len(raw))
		_, err := base64.StdEncoding.Decode(value, []byte(raw))
		return value, err
	case types.Sfixed32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case types.Sfixed64:
		value, err := strconv.ParseInt(raw, 10, 64)
		return value, err
	case types.Sint32:
		value, err := strconv.ParseInt(raw, 10, 32)
		return int32(value), err
	case types.Sint64:
		value, err := strconv.ParseInt(raw, 10, 64)
		return value, err
	}

	return "", fmt.Errorf("unkown type %+v", typed)
}

// StringEmpty returns the given value as a string or a empty string if the value is nil
func StringEmpty(val interface{}) string {
	if val == nil {
		return ""
	}

	return val.(string)
}

// BoolEmpty returns the given value as a bool or a empty bool if the value is nil
func BoolEmpty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatBool(val.(bool))
}

// Int32Empty returns the given value as a int32 or a empty int32 if the value is nil
func Int32Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatInt(int64(val.(int32)), 10)
}

// Uint32Empty returns the given value as a uint32 or a empty uint32 if the value is nil
func Uint32Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatUint(uint64(val.(uint32)), 10)
}

// Int64Empty returns the given value as a int64 or a empty int64 if the value is nil
func Int64Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatInt(val.(int64), 10)
}

// Uint64Empty returns the given value as a uint64 or a empty uint64 if the value is nil
func Uint64Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatUint(val.(uint64), 10)
}

// Float64Empty returns the given value as a float64 or a empty float64 if the value is nil
func Float64Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatFloat(val.(float64), 'E', -1, 64)
}

// Float32Empty returns the given value as a float32 or a empty float32 if the value is nil
func Float32Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return strconv.FormatFloat(float64(val.(float32)), 'E', -1, 32)
}

// BytesBase64Empty returns the given bytes buffer as a base64 string or a empty string if the value is nil
func BytesBase64Empty(val interface{}) string {
	if val == nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(val.([]byte))
}
