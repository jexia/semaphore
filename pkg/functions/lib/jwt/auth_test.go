package jwt

import "testing"

func TestGetAuthorizartionValue(t *testing.T) {
	type test struct {
		value    interface{}
		expected string
		error    error
	}

	var tests = map[string]test{
		"should return an error when the value is not a string": {
			value: 42,
			error: errInvalidValueType,
		},
		"should return an error when the value is invalid": {
			value: "invalid",
			error: errMalformedAuthValue,
		},
		"should return an error when method is not supported": {
			value: "unknown method",
			error: errUnsupportedAuthMethod{kind: "unknown"},
		},
		"should return a JWT when auth value is valid": {
			value:    "Bearer valid.jwt.token",
			expected: "valid.jwt.token",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := getAuthorizartionValue(test.value)

			if err != test.error {
				t.Errorf("error (%v) was expected to be (%v)", err, test.error)
			}

			if actual != test.expected {
				t.Errorf("output (%s) was expected to be (%s)", actual, test.expected)
			}
		})
	}
}
