package sprintf

import "github.com/jexia/semaphore/cmd/semaphore/functions/strings/sprintf/formatter"

var defaultScanner Scanner

func init() {
	var (
		formatters = []Formatter{
			formatter.String{},
			formatter.JSON{},
		}

		detector = NewRadix()
	)

	for _, formatter := range formatters {
		if err := detector.Register(formatter); err != nil {
			panic(err)
		}
	}

	defaultScanner = NewDefaultScanner(detector)
}
