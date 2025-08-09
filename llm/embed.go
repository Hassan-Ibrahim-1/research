package llm

import (
	"slices"

	"github.com/Hassan-Ibrahim-1/research/command"
)

func embed(str []byte, rng command.Range, data []byte) []byte {
	s := slices.Delete(str, rng.Start, rng.End)
	s = slices.Insert(s, rng.Start, data...)
	return s
}
