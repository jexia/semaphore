package http

// func TestParseListenerOptions(t *testing.T) {
// 	duration := time.Second
// 	options := specs.Options{
// 		ReadTimeoutOption:  duration.String(),
// 		WriteTimeoutOption: duration.String(),
// 	}
//
// 	result, err := ParseListenerOptions(options)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	if result.ReadTimeout != duration {
// 		t.Fatalf("unexpected read timeout %+v, expected %+v", result.ReadTimeout, duration)
// 	}
//
// 	if result.WriteTimeout != duration {
// 		t.Fatalf("unexpected write timeout %+v, expected %+v", result.ReadTimeout, duration)
// 	}
// }
