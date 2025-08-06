package prompt

import (
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

var (
	blockCursorStyle = lg.NewStyle().
				Foreground(lg.Color("0")).
				Background(lg.Color("15"))
	blankCursor = []rune(blockCursorStyle.Render(" "))
)

type line struct {
	// runes must never contain a '\n'.
	runes []rune

	// pos must always be in the range 0..=len(runes).
	// pos points to where the next character will be inserted.
	// if it is at the end of the line (pos = len(runes)) then a character will be appended to the line.
	pos int
}

func newLine(maxWidth int) line {
	return line{
		runes: make([]rune, 0, maxWidth),
		pos:   0,
	}
}

func (l *line) addRunes(runes []rune, i int) {
	if slices.Contains(runes, '\n') {
		panic("line.runes cannot have a newline!")
	}

	l.pos += len(runes)

	if i == len(l.runes) {
		l.runes = append(l.runes, runes...)
	} else {
		l.runes = slices.Insert(l.runes, i, runes...)
	}
}

type Model struct {
	viewport viewport.Model
	focused  bool

	lines       []line
	currentLine int
	maxWidth    int

	Cursor cursor.Model

	promptPrefix   string
	characterLimit int
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
		maxWidth:       width - 3,
		// height is meant to be growable
		// what happens if it can't grow anymore?
		// it should page right?
		viewport: vp,
		lines:    make([]line, 0),
	}
}

func (m *Model) addLine() *line {
	log.Println("adding new line")
	m.lines = append(m.lines, newLine(m.maxWidth))
	return &m.lines[len(m.lines)-1]
}

func (m *Model) lineAt(i int) *line {
	if i == 0 && len(m.lines) == 0 {
		log.Println("lineAt adding new line")
		return m.addLine()
	}
	return &m.lines[i]
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

		case tea.KeySpace:
			m.addContent([]rune{' '})
		case tea.KeyEnter:
			m.addLine()
			m.currentLine++
			m.redraw()

		case tea.KeyRight, tea.KeyLeft, tea.KeyUp, tea.KeyDown:
			m.handleArrowKeys(msg)

		default:
			// a control character is sent
			if k := msg.String(); k == "alt+enter" {
				// log.Println("Prompt entered, content:", m.lines)
				cmd = newPromptEnteredMsg(m.String())

				// clear content when the prompt is entered
				m.clear()
			}
		}
	}

	return m, cmd
}

func (m *Model) handleArrowKeys(key tea.KeyMsg) {
	ln := m.lineAt(m.currentLine)

	switch key.Type {
	case tea.KeyRight:
		if ln.pos < len(ln.runes) {
			ln.pos++
			if ln.pos == 49 {
				panic("49 at handleArrowKeys")
			}
		} else if m.currentLine < len(m.lines)-1 {
			m.currentLine++
			m.lineAt(m.currentLine).pos = 0
		}
	case tea.KeyLeft:
		if ln.pos > 0 {
			ln.pos--
		} else if m.currentLine > 0 {
			m.currentLine--
		}
	}

	px, py := m.cursorPos()
	log.Printf("arrow key pos: (%d, %d)\n", px, py)

	m.redraw()
}

func (m *Model) cursorPos() (x, y int) {
	return m.lineAt(m.currentLine).pos, m.currentLine
}

func (m *Model) redraw() {
	m.viewport.SetContent(m.String())
}

func (m *Model) String() string {
	b := strings.Builder{}
	for i, ln := range m.lines {
		runes := ln.runes

		if i == m.currentLine {
			runes = slices.Clone(ln.runes)

			if ln.pos == len(runes) {
				// render the cursor at the end of the line
				runes = append(runes, blankCursor...)
			} else {
				// render at ln.pos
				if ln.pos >= len(runes) {
					log.Println("whoops")
				}
				styled := []rune(blockCursorStyle.Render(string(runes[ln.pos])))
				runes[ln.pos] = styled[0]
				runes = slices.Insert(runes, ln.pos+1, styled[1:]...)
			}

		}

		b.WriteString(string(runes) + "\n")
	}
	return b.String()
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m *Model) clear() {
	m.lines = nil
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Focused() bool {
	return m.focused
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

// possibly returns a PromptResizeMsg
func (m *Model) addContent(runes []rune) tea.Cmd {
	m.writeRunes(runes)
	m.redraw()

	px, py := m.cursorPos()
	log.Printf(
		"cursor pos: (%d, %d)\n",
		px, py,
	)
	return nil
}

func (m *Model) writeRunes(runes []rune) {
	ln := m.lineAt(m.currentLine)
	ln.addRunes(runes, ln.pos)
	m.adjustLines()
	log.Println("brea")
}

func (m *Model) removeChar() {
	ln := m.lineAt(m.currentLine)

	// deferred because of a possible early return
	defer func() {
		m.adjustLines()
		m.redraw()
	}()

	if ln.pos == 0 {
		if m.currentLine > 0 {
			// move n characters from currentLine to the line that we are about to move to.
			// n is the amount of characters is how many characters the above line can hold
			// the above lines cursor is set to the end
			m.currentLine--
			ln = m.lineAt(m.currentLine)
			log.Println("m.currentLine is bigger than 0:", m.currentLine)
		}
		return
	}

	ln.pos--
	ln.runes = slices.Delete(ln.runes, ln.pos, ln.pos+1)
}

func (m *Model) adjustLines() {
	// not using for i := range l.data because we need len(l.data) to be
	// evaluated each iteration because new lines can be appended
	// to the array
	for i := 0; i < len(m.lines); i++ {
		ln := m.lineAt(i)

		if len(ln.runes) >= m.maxWidth {
			// Move the extra characters to the start of the nextLine.
			overflown := ln.runes[m.maxWidth:]
			ln.runes = ln.runes[:m.maxWidth]

			var nextLine *line
			if i < len(m.lines)-1 {
				log.Println("a new line already exists")
				nextLine = m.lineAt(i + 1)

				if ln.pos >= len(ln.runes) && m.currentLine == i {
					// if inserting at the end move the cursor down
					m.currentLine++
					nextLine.pos = 0
				}

			} else {
				nextLine = m.addLine()
				if len(ln.runes) == ln.pos {
					m.currentLine++
				}
			}

			// because we're removing extra characters set ln.pos to the new len
			if ln.pos >= len(ln.runes) {
				ln.pos = len(ln.runes)
			}

			nextLine.addRunes(overflown, 0)
		}
	}
	lastLine := m.lineAt(len(m.lines) - 1)
	if len(lastLine.runes) == 0 && m.currentLine != len(m.lines)-1 {
		m.lines = slices.Delete(m.lines, len(m.lines)-1, len(m.lines))
	}
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

func newPromptResizeMsg(height int) tea.Cmd {
	return func() tea.Msg {
		return PromptResizeMsg{height}
	}
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
