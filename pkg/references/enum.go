package references

import "strconv"

// EnumVal represents a enum value
type EnumVal struct {
	key string
	pos int32
}

// Key returns the key of the given enum value
func (val *EnumVal) Key() string {
	return val.key
}

// Pos returns the enum position of the given enum value
func (val *EnumVal) Pos() int32 {
	return val.pos
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
