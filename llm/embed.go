package llm

import (
	"slices"

	"github.com/Hassan-Ibrahim-1/research/command"
)

// rng is the substring that is replaced by data
func embed(str []byte, rng command.Range, data []byte) []byte {
	if rng.Start == rng.End {
		return str
	}

	s := slices.Delete(str, rng.Start, rng.End)
	s = slices.Insert(s, rng.Start, data...)
	return s
}
