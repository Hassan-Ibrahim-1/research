package command

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestParseCommands(t *testing.T) {
	tests := []struct {
		input    string
		expected []Command
	}{
		{
			"no commands",
			[]Command{},
		},
		{
			"\\@attach-file(file.txt)",
			[]Command{},
		},
		{
			"@attach-file(file.txt) \\@attach-link(example.com)",
			[]Command{NewCommand("attach-file", []string{"file.txt"}, 0, 22)},
		},
		{
			"file: @attach-file(file.txt) link: @attach-link(example.com)",
			[]Command{
				NewCommand("attach-file", []string{"file.txt"}, 6, 28),
				NewCommand("attach-link", []string{"example.com"}, 35, 60),
			},
		},
		{
			"file: @attach-file(file.txt)\nlink: @attach-link(example.com)",
			[]Command{
				NewCommand("attach-file", []string{"file.txt"}, 6, 28),
				NewCommand("attach-link", []string{"example.com"}, 35, 60),
			},
		},
		{
			"file: @attach-file(file.txt, image.png)\nlink: @attach-link(example.com)",
			[]Command{
				NewCommand(
					"attach-file",
					[]string{"file.txt", "image.png"},
					6,
					39,
				),
				NewCommand("attach-link", []string{"example.com"}, 46, 71),
			},
		},
		{
			"file: @attach-file(file.txt, image.png)\nlink: @attach-link(example.com, google.com)",
			[]Command{
				NewCommand(
					"attach-file",
					[]string{"file.txt", "image.png"},
					6,
					39,
				),
				NewCommand(
					"attach-link",
					[]string{"example.com", "google.com"},
					46,
					83,
				),
			},
		},
		{
			`
Hey here's some files @attach-file(image.png, main.go, test.c)
And here's some urls @attach-link(example.com, wikipedia.com)
and here's a markdown file @attach-markdown(input.md).
`,
			[]Command{
				NewCommand(
					"attach-file",
					[]string{"image.png", "main.go", "test.c"},
					23,
					63,
				),
				NewCommand(
					"attach-link",
					[]string{"example.com", "wikipedia.com"},
					85,
					125,
				),
				NewCommand(
					"attach-markdown",
					[]string{"input.md"},
					153,
					179,
				),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cmds := ParseCommands([]byte(tt.input))
			if len(cmds) != len(tt.expected) {
				t.Fatalf(
					"unequal number of commands: got=%d. expected=%d. commands=%s. expected=%s",
					len(cmds),
					len(tt.expected),
					commandSliceString(cmds),
					commandSliceString(tt.expected),
				)
			}
			for i := range cmds {
				_ = testCommandEqual(t, cmds[i], tt.expected[i])
			}
		})
	}
}

func TestParseCommandsString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"no commands", []string{}},
		{"\\@attach-file(file.txt)", []string{}},
		{
			"@attach-file(file.txt) \\@attach-link(example.com)",
			[]string{"@attach-file(file.txt)"},
		},
		{
			"file: @attach-file(file.txt) link: @attach-link(example.com)",
			[]string{"@attach-file(file.txt)", "@attach-link(example.com)"},
		},
		{
			"file: @attach-file(file.txt)\nlink: @attach-link(example.com)",
			[]string{"@attach-file(file.txt)", "@attach-link(example.com)"},
		},
		{
			"file: @attach-file(file.txt, image.png)\nlink: @attach-link(example.com, google.com)",
			[]string{
				"@attach-file(file.txt, image.png)",
				"@attach-link(example.com, google.com)",
			},
		},
		{
			`
Hey here's some files @attach-file(image.png, main.go, test.c)
And here's some urls @attach-link(example.com, wikipedia.com)
and here's a markdown file @attach-markdown(input.md).
`,
			[]string{
				"@attach-file(image.png, main.go, test.c)",
				"@attach-link(example.com, wikipedia.com)",
				"@attach-markdown(input.md)",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cmds := ParseCommands([]byte(tt.input))
			if len(cmds) != len(tt.expected) {
				t.Fatalf(
					"unequal number of commands: got=%d. expected=%d. commands=%s. expected=%+v",
					len(cmds),
					len(tt.expected),
					commandSliceString(cmds),
					tt.expected,
				)
			}
			for i, cmd := range cmds {
				if s := cmd.String(); s != tt.expected[i] {
					t.Errorf(
						"command not equal. got=%q. expected=%q.",
						s,
						tt.expected[i],
					)
				}
			}
		})
	}
}

func TestParseCommandsSubstring(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"no command"},
		{"@attach-file(file.txt)"},
		{"\\@attach-file(file.txt)"},
		{"@attach-file(file.txt, image.png)"},
		{
			"file: @attach-file(file.txt) link: @attach-link(example.com)",
		},
		{
			"@attach-file(file.txt) \\@attach-link(example.com)",
		},
		{
			"file: @attach-file(file.txt)\nlink: @attach-link(example.com)",
		},
		{
			"@attach-file(file.txt)@attach-link(example.com)",
		},
		{
			"file: @attach-file(file.txt)\nlink: @attach-link(example.com)",
		},
		{
			`
Hey here's some files @attach-file(image.png, main.go, test.c)
And here's some urls @attach-link(example.com, wikipedia.com)
and here's a markdown file @attach-markdown(input.md).
`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cmds := ParseCommands([]byte(tt.input))
			for _, cmd := range cmds {
				substr := tt.input[cmd.Loc.Start:cmd.Loc.End]
				if s := cmd.String(); s != substr {
					t.Errorf(
						"Command string not equal to substring. got=%q. expected=%q",
						s,
						substr,
					)
				}
			}
		})
	}
}

func commandSliceString(cmds []Command) string {
	cmdStrings := make([]string, len(cmds))
	for i, cmd := range cmds {
		cmdStrings[i] = cmd.String()
	}

	return "[\n" + strings.Join(cmdStrings, ",\n") + "\n"
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected Command
	}{
		{
			"@attach-file(file.txt)",
			NewCommand("attach-file", []string{"file.txt"}, 0, 22),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cmd, _, err := parseCommand([]byte(tt.input), 0)
			if err != nil {
				t.Fatalf("Failed to parsed command: %v", err)
			}
			_ = testCommandEqual(t, cmd, tt.expected)
		})
	}
}

func TestCommandString(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"@attach-file(file.txt, image.png)"},
		{"@attach-link(example.com, google.com)"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cmd, n, err := parseCommand([]byte(tt.input), 0)
			if err != nil {
				t.Fatalf("Failed to parsed command: %v", err)
			}
			if n != len(tt.input) {
				t.Fatalf(
					"Wrong bytes read. expected=%d. got=%d",
					len(tt.input),
					n,
				)
			}

			if s := cmd.String(); s != tt.input {
				t.Errorf(
					"cmd strings not equal. expected=%q. got=%q",
					tt.input,
					s,
				)
			}
		})
	}
}

func testCommandEqual(t *testing.T, cmd Command, expected Command) bool {
	if cmd.Name != expected.Name {
		t.Errorf(
			"unexpected command name. got=%q. expected=%q",
			cmd.Name,
			expected.Name,
		)
		return false
	}
	if !slices.Equal(cmd.Arguments, expected.Arguments) {
		t.Errorf(
			"unexpected command arguments. got=%q. expected=%q",
			strings.Join(cmd.Arguments, ","),
			strings.Join(expected.Arguments, ","),
		)
		return false
	}
	if cmd.Loc != expected.Loc {
		t.Errorf(
			"unequal Loc. got=%s. expected=%s",
			cmd.Loc.String(),
			expected.Loc.String(),
		)
	}
	return true
}
