package prompt

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	content  string
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

func New(width int) Model {
	if width < 0 {
		panic("width must not be negative")
	}

	vp := viewport.New(width, 1)
	// vp.YPosition = ypos

	return Model{
		promptPrefix:   "> ",
		characterLimit: 1024,
		// height is meant to be growable
		// what happens if it can't grow anymore?
		// it should page right?
		viewport: vp,
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
			// TODO: check if the line is too long. > m.viewport.Width
			m.content += string(msg.Runes)
			m.viewport.SetContent(m.content)
			log.Println("content:", m.content)

		case tea.KeyCtrlC:
			// TODO: TEMPORARY
			return m, tea.Quit

		default:
			// a control character is sent
			// TODO: handle the case where a newline / whitespace is sent
			if k := msg.String(); k == "alt+enter" {
				log.Println("Prompt entered, content:", m.content)
				cmd = newPromptEnteredMsg(m.content)
				// clear content when the prompt is entered
				m.content = ""
			}
		}
	}

	return m, cmd
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
