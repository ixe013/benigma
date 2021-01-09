package benigma

import (
	"testing"

	"github.com/emedvedev/enigma"
)

func TestEncode(t *testing.T) {
    type singleTest struct {
        Model string
        Keyboard string
        Lights string
    }

    tests := []singleTest{
        { "I",  "THISISATEST", "ZPJJSVSPGBW" },
        { "M3", "THISISATEST", "ZPJJSVSPGBW" },
        { "M4", "THISISATEST", "ZPJJSVSPGBW" },
    }

    for i, test_case := range tests {
        machine := builtinModels[test_case.Model]

        encoded := machine.EncodeString(test_case.Keyboard)

        if encoded != test_case.Lights {
            t.Errorf("Run %d: Expected %s but got %s", i, test_case.Keyboard, encoded)
        }
    }
}

func TestSanitize(t *testing.T) {
    type sanitizeTest struct {
        Raw string
        Sanitized string
    }

    tests := []sanitizeTest{
        { "Hello World", "HELLOXWORLD" },
        { "guillaume@paralint.com", "GUILLAUMEPARALINTCOM" },
    }

    for i, test_case := range tests {
        sanitized := enigma.SanitizePlaintext(test_case.Raw)
        if test_case.Sanitized != sanitized {
            t.Errorf("Run %d: Sanitization of %s should yield %s but got %s", i, test_case.Raw, test_case.Sanitized, sanitized)
        }
    }
}
