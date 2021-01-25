package benigma

import (
	"github.com/emedvedev/enigma"
)

/*
Enigma configuration:
  Rotors: [III II I]
  Rotor positions: [A A A]
  Rings: [1 1 1]
  Plugboard: empty
  Reflector: B
*/

var builtinModels = map[string]*enigma.Enigma{
	"I": enigma.NewEnigma(
		bundleRotors( //Rotors
			[]string{"III", "II", "I"},
			[]int{1, 1, 1},
			[]string{"A", "A", "A"},
		),
		"B",        //Reflector
		[]string{}, //Plugboard
	),
	"M3": enigma.NewEnigma(
		bundleRotors( //Rotors
			[]string{"III", "II", "I"},
			[]int{1, 1, 1},
			[]string{"A", "A", "A"},
		),
		"B",        //Reflector
		[]string{}, //Plugboard
	),
	"M4": enigma.NewEnigma(
		bundleRotors( //Rotors
			[]string{"Beta", "III", "II", "I"},
			[]int{1, 1, 1, 1},
			[]string{"A", "A", "A", "A"},
		),
		"B-thin",   //Reflector
		[]string{}, //Plugboard
	),
	"IXE013": enigma.NewEnigma(
		bundleRotors( //rotors
			[]string{"II", "Beta", "VII"},
			[]int{1, 2, 3},
			[]string{"C", "L", "E"},
		),
		"B", //Reflector
		[]string{"AB", "CD", "EF", "GH", "IJ", "KL", "MN", "OP", "QR", "ST", "UV", "WX", "YZ"}, //Plugboard
	),
}

func builtinModelNames() (names []string) {
	for k := range builtinModels {
		names = append(names, k)
	}
	return names
}
