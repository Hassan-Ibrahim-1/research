package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/Hassan-Ibrahim-1/research/llm"
	"github.com/Hassan-Ibrahim-1/research/ui/prompt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	titleStyle = func() lg.Style {
		b := lg.RoundedBorder()
		return lg.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	chatStyle = lg.NewStyle().Align(lg.Center, lg.Center)

	infoStyle = func() lg.Style {
		b := lg.RoundedBorder()
		return titleStyle.BorderStyle(b)
	}()

	errorStyle = lg.NewStyle().Foreground(lg.Color("31")).Bold(true)
)

const glamourStyle = "dark"

type llmResponseStartedMsg struct {
	ch <-chan string
}

type llmPartialResponseMsg struct {
	content string
	ch      <-chan string
}

type llmResponseDoneMsg struct{}

type Model struct {
	viewport viewport.Model
	ready    bool
	prompt   prompt.Model

	// messages between the user and llm that are rendered using glamour
	// when readingLlmResponse is false messages is displayed to the user
	messages string

	// if readingLlmResponse is true this is set to messages + whatever
	// content is being streamed in
	// when readingLlmResponse is true this messages + currentMessage
	// is displayed to the user
	currentMessage *string

	readingLlmResponse bool

	session *llm.Session
}

func New(session *llm.Session) Model {
	return Model{
		session: session,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) redrawViewport(content string) error {
	m.viewport.SetContent(content)
	return nil
}

func (m *Model) onPromptEntered(prompt string) (tea.Cmd, error) {
	r, err := glamour.Render("User: "+prompt+"\n", glamourStyle)
	if err != nil {
		return nil, err
	}
	m.messages += r

	m.prompt.Blur()

	ch, err := m.session.SendPrompt(prompt)
	if err != nil {
		return nil, err
	}

	return func() tea.Msg {
		return llmResponseStartedMsg{ch}
	}, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.prompt.Focus()
		case "esc":
			m.prompt.Blur()
		}

	case tea.WindowSizeMsg:
		// TODO: handle promptView resizes
		m.onWindowResize(msg)

	case prompt.PromptEnteredMsg:
		cmd, err := m.onPromptEntered(msg.Content)
		if err != nil {
			m.reportError(err)
		}
		cmds = append(cmds, cmd)
		m.redrawViewport(m.messages)

	case llmResponseStartedMsg:
		m.startReadingLlmResponse()
		cmds = append(cmds, readResponse(msg.ch))

	case llmPartialResponseMsg:
		if !m.readingLlmResponse {
			panic("impossible state: m.readingLlmResponse must be set to true for llmPartialResponseMsg to be sent")
		}
		if m.currentMessage == nil {
			panic("impossible state: m.currentMessage must not be nil for llmPartialResponseMsg to be sent")
		}

		cmds = append(cmds, readResponse(msg.ch))
		*m.currentMessage += msg.content

		// doing word wrapping here because sometimes the text can get too
		// long for the screen
		wrapped := wordwrap.String(m.messages+*m.currentMessage, m.viewport.Width-3)
		m.redrawViewport(wrapped)

	case llmResponseDoneMsg:
		if !m.readingLlmResponse {
			panic("impossible state: m.readingLlmResponse must be set to true for llmResponseDoneMsg to be sent")
		}
		if m.currentMessage == nil {
			panic("impossible state: m.currentMessage must not be nil for llmResponseDoneMsg to be sent")
		}
		r, err := glamour.Render(*m.currentMessage, glamourStyle)
		if err != nil {
			m.reportError(err)
		} else {
			m.messages += r
		}
		m.stopReadingLlmResponse()
		m.redrawViewport(m.messages)
	}

	if !m.prompt.Focused() {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) startReadingLlmResponse() {
	m.readingLlmResponse = true
	m.currentMessage = new(string)
	m.prompt.SetCanEnterMessage(false)
}

func (m *Model) stopReadingLlmResponse() {
	m.currentMessage = nil
	m.readingLlmResponse = false
	m.prompt.SetCanEnterMessage(true)
}

func readResponse(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return llmResponseDoneMsg{}
		}
		return llmPartialResponseMsg{content: msg, ch: ch}

	}
}

func (m *Model) reportError(err error) {
	log.Println("err:", err)
	m.messages += "err"
	m.redrawViewport(m.messages)
}

func (m Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.headerView(),
		m.chatView(),
		// m.footerView(),
		m.promptView(),
	)
}

func (m *Model) onWindowResize(ws tea.WindowSizeMsg) {
	headerHeight := lg.Height(m.headerView())
	footerHeight := lg.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	viewportWidth := ws.Width

	if !m.ready {
		m.prompt = prompt.New(viewportWidth, 10)

		viewportHeight :=
			ws.Height - (verticalMarginHeight + lg.Height(m.promptView()))

		m.viewport = viewport.New(viewportWidth, viewportHeight)
		m.viewport.YPosition = headerHeight

		// m.prompt.SetYPosition(viewportHeight + footerHeight)
		m.prompt.Focus()
		m.prompt.SetStyle(lg.NewStyle().BorderStyle(lg.RoundedBorder()))
		m.prompt.SetCanEnterMessage(true)

		m.viewport.Style = lg.NewStyle().BorderStyle(lg.RoundedBorder())
		// BorderForeground(lg.Color("62")).
		// Padding(2).Margin(10)

		m.ready = true
	} else {
		viewportHeight :=
			ws.Height - (verticalMarginHeight + lg.Height(m.promptView()))

		m.viewport.Height = viewportHeight
		m.viewport.Width = viewportWidth
	}
}

func (m *Model) headerView() string {
	title := titleStyle.Render("Research")
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(title)))
	return lg.JoinHorizontal(lg.Center, title, line)
}

func (m *Model) chatView() string {
	return m.viewport.View()
}

func (m *Model) footerView() string {
	info := infoStyle.Render(
		fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100),
	)
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(info)))
	return lg.JoinHorizontal(lg.Center, info, line)
}

func (m *Model) promptView() string {
	return m.prompt.View()
}
