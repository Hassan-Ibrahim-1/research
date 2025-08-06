package prompt

import "slices"

// returns -1 if target was not found
// starts from the end of the string and goes to the beginning
func runeIndexReverse(runes []rune, target rune) int {
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == target {
			return i
		}
	}
	return -1
}

// returns -1 if target was not found
func runeIndex(runes []rune, target rune) int {
	for i, r := range runes {
		if r == target {
			return i
		}
	}
	return -1
}

func runesEndsWith(runes []rune, target []rune) bool {
	if len(runes) == 0 || len(target) > len(runes) {
		return false
	}
	endSlice := runes[len(runes)-len(target):]
	return slices.Equal(endSlice, target)
}
