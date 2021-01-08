package benigma

import (
	"testing"
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
