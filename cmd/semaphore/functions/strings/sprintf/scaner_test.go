package sprintf

import (
	"log"
	"testing"

	"github.com/jexia/semaphore/cmd/semaphore/functions/strings/sprintf/formatter"
)

func TestDefaultScanner(t *testing.T) {
	var (
		formatters = []Formatter{
			formatter.String{},
			formatter.JSON{},
		}

		detector = NewRadix()
	)

	for _, formatter := range formatters {
		if err := detector.Register(formatter); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	var scanner = NewDefaultScanner(detector)

	input := "%json: what is that %s, does anybody know?"

	tokens, err := scanner.Scan(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	log.Println(tokens)
}
