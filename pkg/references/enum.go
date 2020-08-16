package references

import "strconv"

// EnumVal represents a enum value
type EnumVal struct {
	key string
	pos int32
}

// MarshalJSON custom marshal implementation mainly used for testing purposes
func (val *EnumVal) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(val.key)), nil
}

// UnmarshalJSON custom unmarshal implementation mainly used for testing purposes
func (val *EnumVal) UnmarshalJSON([]byte) error {
	return nil
}

// Enum value type
func Enum(key string, pos int32) *EnumVal {
	return &EnumVal{
		key: key,
		pos: pos,
	}
}
