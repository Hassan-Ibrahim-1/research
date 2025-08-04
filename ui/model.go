package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
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
)

type Model struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func NewModel(content string) *Model {
	return &Model{
		content: content,
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "l":
			m.viewport.Width = max(0, m.viewport.Width-1)
		case "h":
			m.viewport.Height = max(0, m.viewport.Height-1)
		}

	case tea.WindowSizeMsg:
		headerHeight := lg.Height(m.headerView())
		footerHeight := lg.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		viewportWidth := msg.Width
		viewportHeight := msg.Height - verticalMarginHeight

		if !m.ready {
			m.viewport = viewport.New(viewportWidth, viewportHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)

			m.viewport.Style = lg.NewStyle().
				BorderStyle(lg.RoundedBorder()).
				BorderForeground(lg.Color("62")).
				Padding(2).Margin(10)

			m.ready = true
		} else {
			m.viewport.Height = viewportHeight
			m.viewport.Width = viewportWidth
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.headerView(),
		m.chatView(),
		m.footerView(),
	)
}

func (m *Model) headerView() string {
	title := titleStyle.Render("Pager")
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
