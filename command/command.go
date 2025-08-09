package command

import (
	"bytes"
	"fmt"
	"strings"
)

// Start..<End
type Range struct {
	Start int
	End   int
}

func (r Range) String() string {
	return fmt.Sprintf("%d..%d", r.Start, r.End)
}

type Command struct {
	Loc Range

	Name      string
	Arguments []string
}

func (c Command) String() string {
	args := strings.Join(c.Arguments, ", ")
	return fmt.Sprintf("@%s(%s)", c.Name, args)
}

func NewCommand(name string, arguments []string, start, end int) Command {
	return Command{
		Name:      name,
		Arguments: arguments,
		Loc:       Range{Start: start, End: end},
	}
}

func ParseCommands(str []byte) []Command {
	cmds := []Command{}

	for i := 0; i < len(str); {
		ch := str[i]
		if ch != '@' || matchPrevious(str, i, '\\') {
			i += 1
			continue
		}

		cmd, n, err := parseCommand(str, i)
		if err != nil {
			i += 1
			continue
		}

		i += n
		cmds = append(cmds, cmd)
	}

	return cmds
}

// n is the amount of characters read for the entire command
func parseCommand(b []byte, start int) (command Command, n int, err error) {
	b = b[start:]
	if b[0] != '@' {
		return Command{}, 0, fmt.Errorf(
			"The first character in the command must be '@' got=%s",
			string(b),
		)
	}
	commandName, err := readUntil(b[1:], '(')
	if err != nil {
		return Command{}, 0, fmt.Errorf("Expected '('")
	}

	nameLen := len(commandName)

	args, err := readUntil(b[nameLen+2:], ')')
	if err != nil {
		return Command{}, 0, fmt.Errorf("Expected ')'")
	}
	commandArgs := parseArguments(args)

	// 3 represents '@' + '(' + ')'
	n = len(args) + nameLen + 3
	return NewCommand(string(commandName), commandArgs, start, n+start), n, nil
}

func parseArguments(args []byte) []string {
	if len(args) == 0 {
		return []string{""}
	}

	split := bytes.Split(args, []byte(","))
	ret := make([]string, len(split))
	for i, b := range split {
		ret[i] = string(bytes.TrimSpace(b))
	}
	return ret
}

// returns an error if the delimter was not found
func readUntil(b []byte, delimiter byte) ([]byte, error) {
	i := bytes.IndexByte(b, delimiter)
	if i == -1 {
		return nil, fmt.Errorf(
			"Delimter %c not found in %s",
			delimiter,
			string(b),
		)
	}
	return b[:i], nil
}

// returns if i-1 is expected
func matchPrevious(b []byte, i int, expected byte) bool {
	if i > 0 {
		return b[i-1] == expected
	}
	return false
}
