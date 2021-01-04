package benigma

import (
	"github.com/emedvedev/enigma"
)

var builtinModels = map[string]enigma.Enigma{
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
	"IXE013": enigma.Enigma{
		Reflector: *enigma.NewReflector("YIZMPKRSHAVCWTBODUQNGJLXFE", "A"),
		Rotors: []*enigma.Rotor{
			enigma.NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", "II", "E"),
			enigma.NewRotor("LEYJVCNIXWPBQMDRTAKZGFUHOS", "Beta", ""),
			enigma.NewRotor("NZJHGRCXMYSWBOUFAIVLPEKQDT", "VII", "ZM"),
		},
	},
}

func builtinModelNames() (names []string) {
	for k := range builtinModels {
		names = append(names, k)
	}
	return names
}
