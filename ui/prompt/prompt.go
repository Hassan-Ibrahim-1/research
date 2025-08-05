package prompt

import (
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	content  lines
	focused  bool

	promptPrefix   string
	characterLimit int

	// TODO: placeholder
}

// PromptEnteredMsg is sent when alt+enter is pressed when the prompt area is focused.
type PromptEnteredMsg struct {
	Content string
}

func newPromptEnteredMsg(content string) tea.Cmd {
	return func() tea.Msg {
		return PromptEnteredMsg{content}
	}
}

// PromptResizeMsg is sent when the prompt shrinks or grows
type PromptResizeMsg struct {
	Height int
}

type line struct {
	runes []rune
	pos   int
}

func newLine(maxWidth int) line {
	return line{
		runes: make([]rune, 0, maxWidth),
		pos:   -1,
	}
}

func newLines(maxWidth int) lines {
	return lines{
		maxWidth: maxWidth,
		data:     []line{newLine(maxWidth)},
	}
}

type lines struct {
	data        []line
	maxWidth    int
	currentLine int
}

func (l *lines) String() string {
	b := strings.Builder{}

	for _, line := range l.data {
		b.WriteString(string(line.runes) + "\n")
	}

	return b.String()
}

func (l *lines) clear() {
	l.data = nil
}

func (l *lines) writeRunes(runes []rune) {
	if len(l.data) == 0 {
		l.addLine()
	}

	ln := &l.data[l.currentLine]

	// this might not be what i want???
	// if there is a line after the currentLine then
	// adjustLines will just merge the two, unless this line has a \n
	// also what if a line is just a newline and i start writing to it.
	// it should still have that newline in the end right?
	if len(ln.runes) > 0 && ln.runes[len(ln.runes)-1] == '\n' {
		ln = l.addLine()
		// TODO: this is bugged
		// l.currentLine++
	}
	// ln.runes = append(ln.runes, runes...)
	ln.addRunes(runes, len(ln.runes))

	previousLen := len(l.data)
	atEnd := l.currentLine == len(l.data)-1

	l.adjustLines()

	if atEnd && len(l.data) > previousLen {
		l.currentLine++
	} else if atEnd && len(l.data) < previousLen {
		l.currentLine--
	}

	// if the current line is full then start writing to the next line
	// but if writing in the middle of a line and the line starts overflowing
	// then move the overflown data to the next line and keep doing this until each line
	// fits within maxWidth. add more lines if needed
	// how would removing characters with backspace work? maybe add another function
	// called removeChar which removes the current character on the current line
	// removeChar will get called when backspace is pressed
	// removeChar should also handle the case where overflow doesn't happen anymore
	// so characters will move back

	// maybe have a separate function that adjusts all lines and makes sure that their
	// newlines should also be handled properly
}

func (l *lines) removeChar() {
	log.Println("deleteing a character")
	ln := &l.data[l.currentLine]
	if len(ln.runes) == 0 {
		if l.currentLine == 0 {
			return
		}
		return
	}

	ln.pos = clamp(ln.pos, 0, len(ln.runes)-1)
	ln.runes = slices.Delete(ln.runes, ln.pos, ln.pos+1)
	ln.pos--
	l.adjustLines()
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (l *lines) adjustLines() {
	if l.maxWidth <= 0 {
		panic("l.maxWidth must be a positive integer")
	}

	// go through each line and if the line is not at max width and doesn't end with a newline
	// then check if a next line exists and if it does move its first n characters to the currentLine
	// where n is maxWidth - len(line.data)
	//
	// if the current line is too long then move the overflown data to the next line if it exists
	// or add a new line and move that data to that line
	//
	// if the currentLine is empty then remove it

	// what if a line contains a new line but still has characters after that newline?

	var linesToRemove []int

	// not using for i := range l.data because we need len(l.data) to be
	// evaluated each iteration because new lines can be appended
	// to the array
	for i := 0; i < len(l.data); i++ {
		ln := &l.data[i]

		var nextLine *line
		if i != len(l.data)-1 {
			nextLine = &l.data[i+1]
		}

		endsWithNewLine := false
		if len(ln.runes) > 0 {
			endsWithNewLine = ln.runes[len(ln.runes)-1] == '\n'
		}

		if !endsWithNewLine && len(ln.runes) < l.maxWidth && nextLine != nil {
			// move the next line's n characters to the currentLine
			nextLine := &l.data[i+1]
			n := min(l.maxWidth-len(ln.runes), len(nextLine.runes))

			// TODO: test for off by 1
			newRunes := nextLine.runes[:n]
			nextLine.runes = nextLine.runes[n:]
			ln.addRunes(newRunes, len(ln.runes))
			// ln.runes = append(ln.runes, newRunes...)
		} else if len(ln.runes) >= l.maxWidth {
			overflown := ln.runes[l.maxWidth:]
			ln.runes = ln.runes[:l.maxWidth]

			if nextLine == nil {
				nextLine = l.addLine()
			}
			nextLine.addRunes(overflown, 0)
		}

		// if a line contains a new line in the middle at index n
		// create a new line of runes[n:]
		// the original line would be runes[:n]
		for {
			n := runeIndex(ln.runes, '\n')
			if n == -1 || n == len(ln.runes)-1 {
				break
			}
			afterNewLine := ln.runes[n+1:]
			ln.runes = ln.runes[:n+1]
			if nextLine == nil {
				nextLine = l.addLine()
			}
			nextLine.addRunes(afterNewLine, 0)
			// nextLine.runes = slices.Insert(nextLine.runes, 0, afterNewLine...)
		}

		if len(ln.runes) == 0 {
			linesToRemove = append(linesToRemove, i)
		}
	}

	for i := len(linesToRemove) - 1; i >= 0; i-- {
		lineToRemove := linesToRemove[i]
		// can't remove the first line ever
		if lineToRemove == 0 {
			continue
		}

		if lineToRemove == l.currentLine {
			l.currentLine--
		}
		l.data = slices.Delete(l.data, lineToRemove, lineToRemove+1)
	}
}

func (l *line) addRunes(runes []rune, i int) {
	l.pos += len(runes)
	log.Println("current line pos:", l.pos)
	l.runes = slices.Insert(l.runes, i, runes...)
}

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

func (l *lines) addLine() *line {
	l.data = append(l.data, newLine(l.maxWidth))
	return &l.data[len(l.data)-1]
}

func newPromptResizeMsg(height int) tea.Cmd {
	return func() tea.Msg {
		return PromptResizeMsg{height}
	}
}

func New(width int) Model {
	if width < 0 {
		panic("width must not be negative")
	}

	vp := viewport.New(width, 20)
	// vp.YPosition = ypos

	return Model{
		promptPrefix:   "> ",
		characterLimit: 1024,
		// height is meant to be growable
		// what happens if it can't grow anymore?
		// it should page right?
		viewport: vp,
		content:  newLines(width),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			m.addContent(msg.Runes)

		case tea.KeyCtrlC:
			// TODO: TEMPORARY
			return m, tea.Quit

		case tea.KeyBackspace:
			m.removeChar()

		default:
			// a control character is sent
			if ch := isWhitespace(msg); ch != nil {
				m.addContent([]rune{*ch})
			} else if k := msg.String(); k == "alt+enter" {
				log.Println("Prompt entered, content:", m.content)
				cmd = newPromptEnteredMsg(m.content.String())

				// clear content when the prompt is entered
				m.content.clear()
			}
		}
	}

	return m, cmd
}

func isWhitespace(msg tea.KeyMsg) *rune {
	switch msg.Type {
	case tea.KeySpace:
		ch := ' '
		return &ch
	case tea.KeyEnter:
		ch := '\n'
		return &ch
	case tea.KeyTab:
		ch := '\t'
		return &ch
	}
	return nil
}

// possibly returns a PromptResizeMsg
func (m *Model) addContent(runes []rune) tea.Cmd {
	m.content.writeRunes(runes)
	m.viewport.SetContent(m.content.String())
	log.Println("content:", m.content)

	return nil
}

func (m *Model) removeChar() {
	m.content.removeChar()
	m.viewport.SetContent(m.content.String())
}

func (m Model) View() string {
	// TODO: draw cursor
	return m.viewport.View()
}

func (m *Model) SetFocus(f bool) {
	m.focused = f
}

func (m *Model) GetFocus() bool {
	return m.focused
}

func (m *Model) ToggleFocus() {
	m.focused = !m.focused
}

func (m *Model) SetStyle(style lg.Style) {
	m.viewport.Style = style
}

func (m *Model) Height() int {
	return m.viewport.Height
}

func (m *Model) SetYPosition(ypos int) {
	m.viewport.YPosition = ypos
}
