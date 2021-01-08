package benigma

import (
    "github.com/emedvedev/enigma"
)

func bundleRotors(rotors []string, rings []int, position []string) []enigma.RotorConfig {
    config := make([]enigma.RotorConfig, len(rotors))

    for index, rotor := range rotors {
        ring := rings[index]
        value := position[index][0]
        config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
    }

    return config
}

