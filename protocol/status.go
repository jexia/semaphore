package protocol

// StatusSuccess checks whether the given status code is a success
func StatusSuccess(code int) bool {
	if code >= 200 && code < 300 {
		return true
	}

	return false
}
