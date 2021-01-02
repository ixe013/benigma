package benigma

import (
	"github.com/emedvedev/enigma"
)

var BuiltinModels = map[string]enigma.Enigma{
	"I": enigma.Enigma{
		Reflector: *enigma.NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT", "B"),
		Rotors: []*enigma.Rotor{
			enigma.NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", "III", "V"),
			enigma.NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
			enigma.NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", "I", "Q"),
		},
	},
	"M3": enigma.Enigma{
		Reflector: *enigma.NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT", "B"),
		Rotors: []*enigma.Rotor{
			enigma.NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", "III", "V"),
			enigma.NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
			enigma.NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", "I", "Q"),
		},
	},
	"M4": enigma.Enigma{
		Reflector: *enigma.NewReflector("ENKQAUYWJICOPBLMDXZVFTHRGS", "B-thin"),
		Rotors: []*enigma.Rotor{
			enigma.NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", "III", "V"),
			enigma.NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
			enigma.NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", "I", "Q"),
		},
	},
}

func BuiltinModelNames() (names []string) {
    for k := range BuiltinModels {
        names = append(names, k)
    }
    return names
}
