package prompt

import (
	"fmt"
	"testing"
)

// TODO: test lines.String()

func TestLines_writeRunes(t *testing.T) {
	type Action struct {
		stringToAdd   string
		lineIndex     int
		expectedLines []string
	}

	tests := []struct {
		l        []string
		maxWidth int
		actions  []Action
	}{
		{
			l:        []string{"Hello"},
			maxWidth: 128,
			actions: []Action{
				{", World", 0, []string{"Hello, World"}},
			},
		},
		{
			l:        []string{"Hello\n"},
			maxWidth: 128,
			actions: []Action{
				{", World", 0, []string{"Hello\n", ", World"}},
			},
		},
		{
			l:        []string{"Hello"},
			maxWidth: 6,
			actions: []Action{
				{", World", 0, []string{"Hello,", " World"}},
			},
		},
		{
			l:        []string{"Hello, World", "foo"},
			maxWidth: 6,
			actions: []Action{
				{"bar", 0, []string{"Hello,", " World", "barfoo"}},
			},
		},
		{
			l:        []string{"Hello, World\n", "foo"},
			maxWidth: 128,
			actions: []Action{
				{"bar", 1, []string{"Hello, World\n", "foobar"}},
			},
		},
		{
			l:        []string{"Hello, World\n", "foo\n"},
			maxWidth: 128,
			actions: []Action{
				{"bar", 1, []string{"Hello, World\n", "foo\n", "bar"}},
			},
		},
		{
			l:        []string{"Hello, World", "foo"},
			maxWidth: 1,
			actions: []Action{
				{
					"bar",
					0,
					[]string{
						"H",
						"e",
						"l",
						"l",
						"o",
						",",
						" ",
						"W",
						"o",
						"r",
						"l",
						"d",
						"b",
						"a",
						"r",
						"f",
						"o",
						"o",
					},
				},
			},
		},
		{
			l:        []string{"Hello, World", "foo"},
			maxWidth: 1,
			actions: []Action{
				{
					"bar",
					1,
					[]string{
						"H",
						"e",
						"l",
						"l",
						"o",
						",",
						" ",
						"W",
						"o",
						"r",
						"l",
						"d",
						"f",
						"o",
						"o",
						"b",
						"a",
						"r",
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			l := lines{
				data:     linesFromStrings(tt.l),
				maxWidth: tt.maxWidth,
			}

			for _, action := range tt.actions {
				l.currentLine = action.lineIndex
				l.writeRunes([]rune(action.stringToAdd))
				testLinesEqual(t, l, action.expectedLines)
			}
		})
	}
}

func TestLines_adjustLines(t *testing.T) {
	tests := []struct {
		lines    []string
		maxWidth int
		expected []string
	}{
		{[]string{"Hello, World"}, 128, []string{"Hello, World"}},
		{[]string{"Hello, World"}, 5, []string{"Hello", ", Wor", "ld"}},
		{[]string{"Hello", ", World"}, 128, []string{"Hello, World"}},
		{[]string{"Hello\n", ", World"}, 128, []string{"Hello\n", ", World"}},
		{
			[]string{"Hello\n", "", ", World"},
			128,
			[]string{"Hello\n", ", World"},
		},
		{
			[]string{"Hello, World\n", "foobar"},
			5,
			[]string{"Hello", ", Wor", "ld\n", "fooba", "r"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			l := lines{
				data:     linesFromStrings(tt.lines),
				maxWidth: tt.maxWidth,
			}
			l.adjustLines()
			testLinesEqual(t, l, tt.expected)
		})
	}
}

func testLinesEqual(t *testing.T, l lines, expectedLines []string) {
	if len(l.data) != len(expectedLines) {
		t.Errorf(
			"Lines length not equal: expected=%d\n got=%d\nLines not equal: got=%+v. expected=%+v",
			len(expectedLines),
			len(l.data),
			linesToStrings(l.data),
			expectedLines,
		)
		return
	}
	for i := range expectedLines {
		expected := expectedLines[i]
		ln := string(l.data[i].runes)
		if ln != expected {
			t.Errorf(
				"Line [%d] not equal: got=%s. expected=%s",
				i+1,
				ln,
				expected,
			)
		}
	}
}

func linesFromStrings(s []string) []line {
	ret := make([]line, len(s))
	for i, ln := range s {
		ret[i] = line{runes: []rune(ln)}
	}
	return ret
}

func linesToStrings(lns []line) []string {
	s := make([]string, len(lns))
	for i, ln := range lns {
		s[i] = string(ln.runes)
	}
	return s
}
